from airflow import DAG
from airflow.operators.python_operator import PythonOperator
from airflow.models import Variable
from airflow.utils.dates import days_ago

from intake.spark_etl import spark_run_etl
from intake.spark_write_to_db import spark_write_to_db


def create_intake_dag(src_name):
    """
    Template Function to Dynamically Define DAG's in Airflow.

    :params src_name (str) - name of source.
        Note: This should correspond to the following:
            1. source .yml config file (intake/sources/co2_weekly_mlo/co2_weekly_mlo_config.yml)
                ex: 'co2_weekly_mlo'
            2. Table Name in the planetpulse PostgreSQL Database
    """
    MASTER_EMAILS = Variable.get('MASTER_EMAIL', default_var=['testemail@gmail.com'])
    DB_USER = Variable.get('DB_USER', default_var=None)
    DB_PSWRD = Variable.get('DB_PSWRD', default_var=None)
    DB_URL = Variable.get('DB_URL', default_var=None)
    PLANET_PULSE_DATA_HOME = Variable.get('PLANET_PULSE_DATA_HOME', default_var="opt/airflow/planet-pulse-data")
    # TODO - we should remove aws env vars and give infra
    # the iam role necessary to access other aws services...
    # For now, a hack will suffice...
    AWS_ACCESS_KEY_ID = Variable.get('AWS_ACCESS_KEY_ID', default_var=None)
    AWS_SECRET_ACCESS_KEY = Variable.get('AWS_SECRET_ACCESS_KEY', default_var=None)

    default_args = {
        'owner': 'tkobil',
        'start_date': days_ago(1),
        'depends_on_past': False,
        'email': MASTER_EMAILS, # sample email for github
        'email_on_failure': False
    }

    with DAG(
        'etl_{}'.format(src_name),
        default_args=default_args,
        description='ETL from NOAA for {}'.format(src_name),
        schedule_interval="@daily",
    ) as dag:
        # TODO - switch PythonOperators to SparkSubmit Operators
        
        date_partition = """{{execution_date.strftime("y=%Y/m=%m/d=%d")}}"""
        output_path = f"{PLANET_PULSE_DATA_HOME}/{src_name}/{date_partition}"

        get_data_task = PythonOperator(
            task_id="get_data_{}".format(src_name),
            python_callable=spark_run_etl,
            op_args=[src_name, output_path, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY]
        )

        table_name = src_name

        write_to_db_task = PythonOperator(
            task_id="write_to_db_{}".format(src_name),
            python_callable=spark_write_to_db,
            op_args=[output_path, table_name, DB_USER, DB_PSWRD, DB_URL, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY]
        )

        get_data_task >> write_to_db_task
    
    return dag
