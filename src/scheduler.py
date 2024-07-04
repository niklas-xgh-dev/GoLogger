from apscheduler.schedulers.blocking import BlockingScheduler
from log_collector import LogCollector
from log_parser import LogParser
from db_handler import DBHandler

class LogScheduler:
    def __init__(self):
        self.scheduler = BlockingScheduler()
        self.log_collector = LogCollector()
        self.log_parser = LogParser()
        self.db_handler = DBHandler()

    def collect_and_store_logs(self):
        logs = self.log_collector.collect_system_logs()
        parsed_logs = self.log_parser.parse_log(logs)
        if parsed_logs:
            self.db_handler.insert_log(parsed_logs)
        else:
            print("No logs to insert")

    def start(self):
        self.scheduler.add_job(self.collect_and_store_logs, 'interval', seconds=30)
        print("Scheduler started. Press Ctrl+C to exit.")
        self.scheduler.start()

if __name__ == "__main__":
    log_scheduler = LogScheduler()
    log_scheduler.start()