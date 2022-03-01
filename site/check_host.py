import argparse
from email.policy import default
import multiprocessing
import requests
import time


def get_report_url(site_url):
    request_str = f'https://check-host.net/check-http?host={site_url}'
    res = requests.get(request_str, timeout=60)

    for line in res.text.splitlines():
        if 'check-report' in line:
            pattern = 'href="'
            beg = line.find(pattern) + len(pattern)
            end = line.find('"', beg + 2)
            report_url = line[beg:end]
            return report_url


def get_avail_by_country(report_url):
    res = requests.get(report_url)
    avail_list = []
    pattern = 'check_displayer.display'
    for line in res.text.splitlines():
        beg = line.find(pattern)
        if beg != -1:
            beg += len(pattern) + 2
            country_id = line[beg: beg + 3]
            avail_list.append((country_id, not "null" in line))
    return avail_list


class Worker(multiprocessing.Process):
    def __init__(self, job_q, results_q):
        super().__init__()
        self._job_q = job_q
        self._results_q = results_q

    def run(self):
        while True:
            site_url = self._job_q.get()
            if site_url is None:
                break

            report_url = get_report_url(site_url)
            time.sleep(10)
            avail_list = get_avail_by_country(report_url)
            print('Done: ' + site_url)
            self._results_q.append((site_url, avail_list))
        

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--input-file', default='urls')
    args = parser.parse_args()

    jobs = []
    job_q = multiprocessing.Queue()
    results_q = multiprocessing.Manager().list()

    for i in range(4):
        p = Worker(job_q, results_q)
        jobs.append(p)
        p.start()

    with open(args.input_file, 'r') as f:
        for line in f.readlines():
            site_url = line.replace('\n', '')
            job_q.put(site_url)
    
    for _ in jobs:
        job_q.put(None)

    for j in jobs:
        j.join()

    with open('./site/site_avail_by_country.csv', 'w') as f:
        header = 'site_url,country_id,avail\n'
        f.write(header)
        for res in results_q:
            site_url, avail_list = res
            for x in avail_list:
                country_id, avail = x
                f.write(','.join([site_url, country_id, str(avail)]) + '\n')