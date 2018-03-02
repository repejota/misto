import sys
import time

import schedule
from faker import Faker

def generate_fakelog():
    fake = Faker()
    # :param nb_words: around how many words the sentence should contain
    logstring = fake.sentence(nb_words=12)
    return logstring

def log_stderr(logstring):
    sys.stderr.write(logstring)

def log_stdout(logstring):
    sys.stdout.write(logstring)

def job():
    logstring = generate_fakelog()
    print(logstring)

def main():
    schedule.every(2).seconds.do(job)
    while True:
        schedule.run_pending()
        time.sleep(1)

if __name__ == "__main__":
    job()
    main()