import sys
import time
import random

from faker import Faker

def generate_fakelog():
    fake = Faker()
    # param nb_words - around how many words the sentence should contain
    logstring = fake.sentence(nb_words=12)
    return logstring

def log_stderr(logstring):
    sys.stderr.write("{}{}".format(logstring, "\n"))
    sys.stderr.flush()

def log_stdout(logstring):
    sys.stdout.write("{}{}".format(logstring, "\n"))
    sys.stdout.flush()

def job():
    logstring = generate_fakelog()
    if random.random() > 0.5:
        logstring = "{} {}".format("[stdout]", logstring)
        log_stdout(logstring)
    else:
        logstring = "{} {}".format("[stderr]", logstring)
        log_stderr(logstring)


def main():
    while True:
        job()
        wait = random.randrange(9)+1
        time.sleep(1)

if __name__ == "__main__":
    main()