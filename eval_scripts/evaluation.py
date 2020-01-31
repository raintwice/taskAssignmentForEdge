#！/usr/bin/env python3
import os
import datetime
import subprocess
import time
import math
import random
import pandas as pd
from multiprocessing import Process

# ratePara times per minutes
def nextTime(ratePara):
    return -math.log(1.0-random.random())/ratePara

def startMaster(dptIndex):
    cmdstr = '(./master -dispatcher %d ) 2>&1 | tee master.log' % (dptIndex)
    subprocess.call(cmdstr, shell=True)

def stopMaster():
    cmdstr = "ps -ef|grep ./master|awk '{print $2}' | xargs kill -9"
    subprocess.call(cmdstr, shell=True)

def startNodes(timeAvg):
    time.sleep(1)
    name = ['bw']
    bw_data = pd.read_csv('bandwidth.csv', names=name)
    if not os.path.exists('./nodeinfo'):
        os.makedirs('./nodeinfo')
    strat_port = 50101
    for i in range(50):
        port = strat_port + i
        cmdstr = '(./node -nport %d -bw %f -pscTimeAvg %d -pscTimeSig %d) 2>&1 | tee ./nodeinfo/node_%d.log &' % (port, bw_data.at[i, 'bw'], timeAvg, timeAvg*6, port)
        print(cmdstr)
        subprocess.call(cmdstr, shell=True)
        intervel = nextTime(5)  #5 tasks per minutes
        time.sleep(intervel*60)

    print("all the nodes are booted up!")

def stopNodes():
    cmdstr = "ps -ef|grep ./node|awk '{print $2}' | xargs kill -9"
    subprocess.call(cmdstr, shell=True)

def startClient(rate, trainNum, evalNum):
    time.sleep(1)
    cmdstr = '(./client -rate %d -ptrainNum %d -evalNum %d) 2>&1 | tee client.log' % (rate, trainNum, evalNum)
    subprocess.call(cmdstr, shell=True)

def SingleTest(dptIndex, timeAvg, rate, trainNum, evalNum):
    pm = Process(target=startMaster, args=(dptIndex,))
    pn = Process(target=startNodes, args=(timeAvg,))
    pc = Process(target=startClient, args=(rate, trainNum, evalNum))

    pm.start()
    pn.start()
    pc.start()
    pc.join()
    stopMaster()
    stopNodes()

    testFileName = 'Result-%d-%d-%d' % (dptIndex, timeAvg, rate)
    if not os.path.exists(testFileName):
        os.makedirs(testFileName)
    cmdstr = 'mv master.log %s' % testFileName
    subprocess.call(cmdstr, shell=True)
    cmdstr = 'mv nodeinfo %s' % testFileName
    subprocess.call(cmdstr, shell=True)
    cmdstr = 'mv client.log %s' % testFileName
    subprocess.call(cmdstr, shell=True)
    cmdstr = 'mv EvalTaskResult.csv %s' % testFileName
    subprocess.call(cmdstr, shell=True)
    cmdstr = 'mv PretrainTaskResult.csv %s' % testFileName
    subprocess.call(cmdstr, shell=True)
    print('Evaluation is done. Dispatcher %d, timeAvg %d, rate %d' % (dptIndex, timeAvg, rate))

if __name__ == "__main__":
    presenceTimeAvgList = [10, 15, 20, 25, 0]
    jobArrivalRateList = [90, 120, 150, 180, 210]
    disIndexList = [1, 2, 3, 4, 5]

    #test stop
    #stopMaster()
    #stopNodes()
    #SingleTest(3, 20, 120, 200, 200)
    '''
    for timeAvg in presenceTimeAvgList:
        for rate in jobArrivalRateList:
            for index in disIndexList:
                SingleTest(index, timeAvg, rate, 5000, 20000)
                time.sleep(2)
    '''
    SingleTest(5, 15, 90, 5000, 20000)
    '''
    timeAvg = 15
    for rate in jobArrivalRateList:
        for index in disIndexList:
            SingleTest(index, timeAvg, rate, 5000, 20000)
            time.sleep(2)

    testFileName = 'Result-job-rate'
    if not os.path.exists(testFileName):
        os.makedirs(testFileName)
    cmdstr = 'mv Result* %s' % testFileName
    subprocess.call(cmdstr, shell=True)

    time.sleep(2)
    rate = 120
    for timeAvg in presenceTimeAvgList:
        for index in disIndexList:
            SingleTest(index, timeAvg, rate, 5000, 20000)
            time.sleep(2)

    testFileName = 'Result-presence-time'
    if not os.path.exists(testFileName):
        os.makedirs(testFileName)
    cmdstr = 'mv Result* %s' % testFileName
    subprocess.call(cmdstr, shell=True)
    '''
    '''
    rateList = [150, 180, 210]
    for rate in rateList:
        for index in disIndexList:
            SingleTest(index, 15, rate, 5000, 20000)
    '''