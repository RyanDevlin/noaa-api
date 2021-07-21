FROM apache/airflow:2.0.1-python3.8

USER root

RUN apt-get update
RUN apt-get install openjdk-11-jdk -y

COPY requirements.txt .
COPY dist/pipeline-1.0-py3-none-any.whl .

RUN python -m pip install --upgrade pip
RUN pip install -r requirements.txt
RUN pip install pipeline-1.0-py3-none-any.whl
