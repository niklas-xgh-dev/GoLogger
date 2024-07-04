import psutil
import json
from datetime import datetime
from pythonjsonlogger import jsonlogger
import logging

class LogCollector:
    def __init__(self):
        self.logger = self.setup_logger()

    def setup_logger(self):
        logger = logging.getLogger()
        logHandler = logging.StreamHandler()
        formatter = jsonlogger.JsonFormatter()
        logHandler.setFormatter(formatter)
        logger.addHandler(logHandler)
        logger.setLevel(logging.INFO)
        return logger

    def collect_system_logs(self):
        cpu_percent = psutil.cpu_percent()
        memory_info = psutil.virtual_memory()
        disk_usage = psutil.disk_usage('/')

        log_data = {
            'timestamp': datetime.now().isoformat(),
            'cpu_percent': cpu_percent,
            'memory_percent': memory_info.percent,
            'disk_percent': disk_usage.percent
        }

        self.logger.info(json.dumps(log_data))
        #print("Collected log data:", log_data)
        return log_data

if __name__ == "__main__":
    collector = LogCollector()
    collector.collect_system_logs()