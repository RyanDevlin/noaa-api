# Running Pipeline in Airflow

## One Time Setup
    mkdir ./logs ./plugins
    echo -e "AIRFLOW_UID=$(id -u)\nAIRFLOW_GID=0" > .env
    docker-compose up airflow-init


## Every Time code is changed
    cd ..
    python setup.py bdist_wheel
    docker build --progress=plain -t airflow:2.0.1 .
    cd airflow
    docker-compose up