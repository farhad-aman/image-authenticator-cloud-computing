import base64
import logging

import pika
import psycopg2
import requests

import settings

logging.basicConfig(level=logging.INFO)  # Set logging level as needed


def face_detection(image_binary):
    api_key = "acc_a74eadebb47ac90"
    api_secret = "5755b628df28d25c6de805b8d063d424"

    params = {
        'return_face_id': 1
    }

    response = requests.post(
        'https://api.imagga.com/v2/faces/detections',
        auth=(api_key, api_secret),
        params=params,
        files={'image': image_binary}
    )

    if response.status_code == 200:
        response_data = response.json()
        confidence = 0
        face_id = ''
        if 'faces' in response_data['result'] and response_data['result']['faces']:
            face_info = response_data['result']['faces'][0]
            confidence = int(face_info['confidence'])
            face_id = face_info.get('face_id', '')
        logging.info(f"Face detection: Confidence - {confidence}, Face ID - {face_id}")
        return confidence, face_id
    else:
        return 0, ''


def face_similarity(image1_face_id, image2_face_id):
    api_key = 'acc_a74eadebb47ac90'
    api_secret = '5755b628df28d25c6de805b8d063d424'

    response = requests.get(
        f'https://api.imagga.com/v2/faces/similarity?face_id={image1_face_id}&second_face_id={image2_face_id}',
        auth=(api_key, api_secret)
    )

    if response.status_code == 200:
        result = response.json().get('result', {})
        score = int(result.get('score', 0))
        logging.info(f"Face similarity: Score - {score}")
        return score
    else:
        return 0


def get_images_from_s3(national: str):
    try:
        image1_url = f"https://{settings.S3_BUCKET}.{settings.S3_ENDPOINT}/{national}-1"
        image1_response = requests.get(image1_url)
        image1_response.raise_for_status()
        image1_binary = base64.b64decode(image1_response.content)

        image2_url = f"https://{settings.S3_BUCKET}.{settings.S3_ENDPOINT}/{national}-2"
        image2_response = requests.get(image2_url)
        image2_response.raise_for_status()
        image2_binary = base64.b64decode(image2_response.content)
        logging.info("Images retrieved from S3 successfully")
        return image1_binary, image2_binary
    except requests.exceptions.HTTPError:
        raise
    except requests.exceptions.RequestException as err:
        print(f"Request error occurred: {err}")
        raise
    except Exception as err:
        print(f"Unexpected error occurred: {err}")
        raise


def init_pg():
    try:
        connection = psycopg2.connect(
            host=settings.DB_HOST,
            port=settings.DB_PORT,
            dbname=settings.DB_NAME,
            user=settings.DB_USER,
            password=settings.DB_PASSWORD
        )
        logging.info("PostgreSQL connection established successfully")
        return connection
    except psycopg2.Error as e:
        print("Error: Unable to connect to the database")
        print(e)
        return None


pg_connection = init_pg()


def get_email(national: str):
    if pg_connection is None:
        return None
    try:
        with pg_connection.cursor() as cursor:
            cursor.execute("SELECT email FROM users WHERE national = %s", (national,))
            email = cursor.fetchone()
            if email:
                return email[0]
            else:
                logging.error(f"Error: No user found with national ID: {national}")
                return None
    except psycopg2.Error as e:
        logging.error(f"Error: Unable to fetch email from the database. {e}")
        return None


def reject_user(national: str, email_address: str):
    if pg_connection is None:
        return False
    try:
        with pg_connection.cursor() as cursor:
            cursor.execute("UPDATE users SET state = 'rejected' WHERE national = %s", (national,))
            pg_connection.commit()
        logging.info(f"User with national ID {national} rejected successfully")
        if email_address:
            send_email("Your Image Authentication Request Rejected.", email_address)
        return True
    except psycopg2.Error as e:
        print("Error: Unable to reject user")
        print(e)
        return False


def accept_user(national: str, email_address: str):
    if pg_connection is None:
        return False
    try:
        with pg_connection.cursor() as cursor:
            cursor.execute("UPDATE users SET state = 'accepted' WHERE national = %s", (national,))
            pg_connection.commit()
        logging.info(f"User with national ID {national} accepted successfully")
        if email_address:
            send_email("Your Image Authentication Request Accepted.", email_address)
        return True
    except psycopg2.Error as e:
        print("Error: Unable to accept user")
        print(e)
        return False


def send_email(message: str, receiver: str):
    domain = "sandboxe347ff0a3f0240239da6f835176d38bf.mailgun.org"
    api_key = "7ba7dec521a2f17974c46c43158668c3-5465e583-dc962af4"
    url = f"https://api.mailgun.net/v3/{domain}/messages"
    try:
        response = requests.post(
            url,
            auth=("api", api_key),
            data={
                "from": f"Image Authenticator <mailgun@{domain}>",
                "to": [receiver],
                "subject": "Image Authentication Result",
                "text": message
            }
        )
        if response.status_code == 200:
            logging.info("Email sent successfully.")
        else:
            logging.error(f"Failed to send email. Status code: {response.status_code}")
    except Exception as e:
        logging.error(f"Error occurred while sending email: {str(e)}")


def callback(ch, method, properties, body):
    national = body.decode('utf-8')
    image1_binary, image2_binary = get_images_from_s3(national)
    image1_confidence, image1_face_id = face_detection(image1_binary)
    image2_confidence, image2_face_id = face_detection(image2_binary)
    email_address = get_email(national)
    if image1_confidence < 50 or image2_confidence < 50:
        if not reject_user(national, email_address):
            logging.error(f"Failed to reject user with national ID: {national}")
    else:
        similarity = face_similarity(image1_face_id, image2_face_id)
        if similarity < 80:
            if not reject_user(national, email_address):
                logging.error(f"Failed to reject user with national ID: {national}")
        else:
            if not accept_user(national, email_address):
                logging.error(f"Failed to accept user with national ID: {national}")


class Listener:

    def __init__(self):
        self.rabbit_connection = pika.BlockingConnection(pika.URLParameters(settings.RABBIT_URL))
        self.rabbit_channel = self.rabbit_connection.channel()
        self.rabbit_channel.basic_consume(queue='national', on_message_callback=callback, auto_ack=True)
        logging.info("RabbitMQ connection established successfully")

    def listen(self):
        logging.info("Listening for messages...")
        self.rabbit_channel.start_consuming()


if __name__ == "__main__":
    listener = Listener()
    listener.listen()
