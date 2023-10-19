import os

RABBIT_URL = os.getenv("RABBIT_URL") or ("amqps://mihcucyt:XE1vceH1adMHJ1lutXATIm6PI79F7aUl@cow.rmq2.cloudamqp.com"
                                         "/mihcucyt")

S3_BUCKET = os.getenv("S3_BUCKET") or "image-authenticator"
S3_REGION = os.getenv("S3_REGION") or "default"
S3_ENDPOINT = os.getenv("S3_ENDPOINT") or "s3.ir-thr-at1.arvanstorage.ir"
S3_ACCESS_KEY = os.getenv("S3_ACCESS_KEY") or "1f587d9d-0f73-44b0-a244-eb65a38b9fb9"
S3_SECRET_KEY = os.getenv("S3_SECRET_KEY") or "ea8f9014e2dba906a1e1d16ced994f0ddb567bee05386ccfb13646e58d41263d"


PG_DSN = os.getenv("PG_DSN") or ("postgres://avnadmin:AVNS_oPGfErq0_Hdn06tezKr@farhad-farhadaman7780.a.aivencloud.com"
                                 ":21646/defaultdb?sslmode=require")

# DB_HOST = os.getenv("DB_HOST") or "pg"
# DB_PORT = os.getenv("DB_PORT") or "5423"
# DB_NAME = os.getenv("DB_NAME") or "image-authenticator"
# DB_USER = os.getenv("DB_USER") or "root"
# DB_PASSWORD = os.getenv("DB_PASSWORD") or "root"
