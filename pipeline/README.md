# Data Ingestion Pipeline

## Overview
This is a data ingestion framework to intake enironment-specific data from the following FTP site: https://gml.noaa.gov/aftp/


## Setup
    # Actiate Virtual Environment
    python3 -m venv venv
    source venv/bin/activate

    # Install Requirements
    pip install -r requirements.txt
    brew install adoptopenjdk11
    brew install adoptopenjdk11

    # Build Package
    python setup.py build
    python setup.py install

## To Run
    # We are still in development at this time. The current pipeline will
    # write a sample file to sample_out/.
    # Please ensure this directory is empty prior to running, otherwise
    # we will fail on spark.write.parquet('sample_out')
    python intake/spark_etl.py --source=co2_weekly_mlo

    # Feel free to inspect the output of the pipeline from the spark-shell following 
    # the below commands
    >>pyspark
    >> df = spark.read.parquet('sample_out')
    >> df.show() # inspect dataframe
    >> df.take(1)[0] # inspect a sample record
    >> df.printSchema() # inspect spark-sql schema for data