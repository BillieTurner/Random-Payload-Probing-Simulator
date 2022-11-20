from collections import Counter, defaultdict
import statistics
import sys

import matplotlib.pyplot as plt

def main():
    fileName = sys.argv[1]
    errCounter = Counter()
    timeDict = defaultdict(list)
    with open(fileName) as f:
        for i, line in enumerate(f):
            if i == 0:
                continue
            size, t, err, _ = line.strip().split(",")
            t = t.strip()
            err = err.strip()

            timeDict[int(size)].append(int(t))
            errCounter[err] += 1

    for k, v in errCounter.items():
        print(k, v)

    rst = {}
    for k in timeDict.keys():
        data = timeDict[k]
        avg = sum(data) / len(data) / 1000
        print(k, avg)
        rst[k] = avg
    stdev = statistics.stdev([item for lst in timeDict.values() for item in lst])
    avg = sum(rst.values()) / len(rst)
    print(f'stdev {stdev / 1000}s')
    print(f'avg {avg}s')

    draw_chart(rst)


def draw_chart(data):
    plt.plot(data.keys(), data.values())

    # plt.yticks([29.5, 30, 30.5])
    plt.xlabel("packet size (byte)")
    plt.ylabel("server response time (second)")
    plt.show()


if __name__ == "__main__":
    main()
