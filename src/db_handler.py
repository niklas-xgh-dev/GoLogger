import os
from sqlalchemy import create_engine, Column, Integer, Float, DateTime
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from dotenv import load_dotenv
from datetime import datetime

load_dotenv('config/.env')

Base = declarative_base()

class SystemLog(Base):
    __tablename__ = 'system_logs'

    id = Column(Integer, primary_key=True)
    timestamp = Column(DateTime)
    cpu_percent = Column(Float)
    memory_percent = Column(Float)
    disk_percent = Column(Float)

class DBHandler:
    def __init__(self):
        db_url = os.getenv('DATABASE_URL')
        if db_url is None:
            raise ValueError("DATABASE_URL environment variable is not set")
        self.engine = create_engine(db_url)
        Base.metadata.create_all(self.engine)
        self.Session = sessionmaker(bind=self.engine)

    def insert_log(self, log_data):
        session = self.Session()
        try:
            new_log = SystemLog(**log_data)
            session.add(new_log)
            session.commit()
            #print(f"Inserted log: {log_data}")
        except Exception as e:
            print(f"Error inserting log: {e}")
            session.rollback()
        finally:
            session.close()

if __name__ == "__main__":
    handler = DBHandler()
    sample_log = {
        'timestamp': datetime.now(),
        'cpu_percent': 10.0,
        'memory_percent': 10.0,
        'disk_percent': 10.0
    }
    handler.insert_log(sample_log)