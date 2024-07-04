import json
from dateutil import parser as date_parser

class LogParser:
    def parse_log(self, log_entry):
        try:
            if isinstance(log_entry, str):
                parsed_log = json.loads(log_entry)
            elif isinstance(log_entry, dict):
                parsed_log = log_entry
            else:
                raise ValueError("Log entry must be a string or dictionary")
            
            parsed_log['timestamp'] = date_parser.parse(parsed_log['timestamp'])
            return parsed_log
        except json.JSONDecodeError:
            print(f"Error parsing log entry: {log_entry}")
            return None
        except KeyError:
            print(f"Missing required fields in log entry: {log_entry}")
            return None
        except ValueError as e:
            print(f"Error: {str(e)}")
            return None

if __name__ == "__main__":
    log_parser = LogParser()
    sample_log = '{"timestamp": "2024-07-04T10:00:00", "cpu_percent": 50.0, "memory_percent": 60.0, "disk_percent": 70.0}'
    parsed = log_parser.parse_log(sample_log)
    print(parsed)