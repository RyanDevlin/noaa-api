import argparse
import os
import pkg_resources
from pyspark.sql import SparkSession
from pyspark import SparkContext, SparkConf


CURR_DIR = os.path.dirname(os.path.realpath(__file__))

def spark_write_to_db(source, table_name, db_user, db_psrwd, db_endpoint, aws_access_key_id, aws_secret_access_key):
    """
    Dumb copy of source data (parquet) into DB table using spark. 
    
    :params source (str) - s3 or local path to parquet file
    :params table_name (str) - name of table in planetpulse postgresql db
    :params db_user (str) - username for db access
    :params db_pswrd (str) - password for db access
    :params db_endpoint (str) - endpoint of planetpulse postgresql db
    """
    # Again, This is a work-around since we are running spark locally on EC2...
    # TODO - this is a similar to spark setup in intake/spark_etl.py - let's refactor to use shared utils code...
    # We should also be able to remove some of this when we move from PythonOperator -> SparkSubmitOperator
    os.environ['PYSPARK_SUBMIT_ARGS'] = '--packages org.apache.hadoop:hadoop-common:3.0.0,org.apache.hadoop:hadoop-aws:3.0.0,org.apache.hadoop:hadoop-client:3.0.0 pyspark-shell'
    os.environ['AWS_ACCESS_KEY_ID'] = aws_access_key_id
    os.environ['AWS_SECRET_ACCESS_KEY'] = aws_secret_access_key
    jars_path = pkg_resources.resource_filename('intake.jars', 'postgresql-42.2.23.jar')
    conf = SparkConf().set('spark.jars', jars_path)
    sc = SparkContext(conf=conf)
    sc.setSystemProperty('com.amazonaws.services.s3.enableV4', 'true')
    hadoopConf = sc._jsc.hadoopConfiguration()
    hadoopConf.set("fs.s3a.awsAccessKeyId", aws_access_key_id)
    hadoopConf.set("fs.s3a.awsSecretAccessKey", aws_secret_access_key)
    hadoopConf.set("fs.s3a.endpoint", "s3.amazonaws.com")
    hadoopConf.set('fs.s3a.impl', 'org.apache.hadoop.fs.s3a.S3AFileSystem')
    spark = SparkSession(sc)

    # Always overwrite data. We are always processing all data from
    # source, rather than just new data. No need to load current db data 
    # and decide what to write.
    mode = 'overwrite'
    properties = {"user": db_user, "password": db_psrwd, "driver": "org.postgresql.Driver"}
    df = spark.read.parquet(source)
    df.write.jdbc(url=db_endpoint, table=table_name, mode=mode, properties=properties)

def get_args():
    parser = argparse.ArgumentParser(description="Spark Write to DB CLI")
    parser.add_argument('--input_path', type=str, dest="input_path", help="s3 path to parquet files", required=True)
    parser.add_argument('--table_name', type=str, dest="table_name", help="table name in PostgresSQL DB plane-pulse", required=True)
    return parser.parse_args()


if __name__ == "__main__":
    db_pswrd = os.environ.get('DB_PSRWD')
    db_user = os.environ.get('DB_USER')
    db_endpoint = os.environ.get('DB_ENDPOINT')
    # spark_write_to_db('sample_out', 'co2_weekly_mlo', db_user, db_pswrd, db_endpoint)
