import sys
import time
import random

import schedule
from faker import Faker

def generate_fakelog():
    fake = Faker()
    # :param nb_words: around how many words the sentence should contain
    logstring = fake.sentence(nb_words=12)
    return logstring

def log_stderr(logstring):
    print(logstring, file=sys.stderr)

def log_stdout(logstring):
    print(logstring)

def job():
    logstring = generate_fakelog()
    if random.random() > 0.5:
        log_stdout(logstring)
    else:
        logstring = "{} {}".format("ERROR:", logstring)
        log_stderr(logstring)

def main():
    schedule.every(2).seconds.do(job)
    while True:
        schedule.run_pending()
        time.sleep(1)

if __name__ == "__main__":
    job()
    main()