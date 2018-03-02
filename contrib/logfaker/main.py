import sys
import time
import random
import logging

import schedule
from faker import Faker

def generate_fakelog():
    fake = Faker()
    # param nb_words - around how many words the sentence should contain
    logstring = fake.sentence(nb_words=12)
    return logstring

def log_stderr(logstring):
    print(logstring, file=sys.stderr)

def log_stdout(logstring):
    print(logstring, file=sys.stdout)

def job():
    logstring = generate_fakelog()
    
    if random.random() > 0.5:
        logstring = "{} {}".format("[stdout]", logstring)
        log_stdout(logstring)
    else:
        logstring = "{} {}".format("[stderr]", logstring)
        log_stderr(logstring)

def main():
    wait = random.randrange(10)+1
    schedule.every(1).seconds.do(job)
    while True:
        schedule.run_pending()
        time.sleep(1)

if __name__ == "__main__":
    job()
    main()