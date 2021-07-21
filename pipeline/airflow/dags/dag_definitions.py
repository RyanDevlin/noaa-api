from airflow import DAG

from dag_templates.intake_dag import create_intake_dag

globals()["co2_weekly_mlo"] =  create_intake_dag("co2_weekly_mlo", "opt/airflow/data_files/sample_out")
