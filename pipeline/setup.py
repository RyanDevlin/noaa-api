from setuptools import setup, find_packages


data_files = [
    (
        'intake/sources/co2_weekly_mlo',
        [
            'intake/sources/co2_weekly_mlo/co2_weekly_mlo_config.yml'
        ]
    ),
    (
        'intake/jars/co2_weekly_mlo',
        [
            'intake/jars/postgresql-42.2.23.jar'
        ]
    )
]

setup(
    name="pipeline",
    version="1.0",
    packages=find_packages(),
    data_files=data_files,
    include_package_data=True,
)