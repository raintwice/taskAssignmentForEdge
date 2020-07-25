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
    cmdstr = "ps -ef|grep ./master|awk '{print $2}' | xargs kill "
    subprocess.call(cmdstr, shell=True)

def startNodes(timeAvg, nodeNum):
    time.sleep(1)
    name = ['bw']
    bw_data = pd.read_csv('bandwidth.csv', names=name)
    if not os.path.exists('./nodeinfo'):
        os.makedirs('./nodeinfo')
    strat_port = 50101
    for i in range(nodeNum):
        port = strat_port + i
        cmdstr = '(./node -nport %d -bw %f -pscTimeAvg %d -pscTimeSig %d) 2>&1 | tee ./nodeinfo/node_%d.log &' % (port, bw_data.at[i%50, 'bw'], timeAvg, timeAvg*6, port)
        print(cmdstr)
        subprocess.call(cmdstr, shell=True)
        intervel = nextTime(5)  #5 nodes per minutes
        time.sleep(intervel*60)

    print("all the nodes are booted up!")

def startNodesOnce(timeAvg, nodeArrivalRate):
    time.sleep(1)
    name = ['bw']
    bw_data = pd.read_csv('bandwidth.csv', names=name)
    if not os.path.exists('./nodeinfo'):
        os.makedirs('./nodeinfo')
    strat_port = 50101
    i = 0
    while True :
        port = strat_port + i
        cmdstr = '(./node -nport %d -bw %f -pscTimeAvg %d -pscTimeSig %d -nodeMode 1) 2>&1 | tee ./nodeinfo/node_%d.log &' % (port, bw_data.at[i%50, 'bw'], timeAvg, timeAvg*6, port)
        print(cmdstr)
        subprocess.call(cmdstr, shell=True)
        intervel = nextTime(nodeArrivalRate)
        time.sleep(intervel*60)
        i = i + 1

def startNodesOnceWithCapa(timeAvg, nodeArrivalRate, capa, mType, gType, startPort):
    if nodeArrivalRate == 0:
        return 
    time.sleep(1)
    name = ['bw']
    bw_data = pd.read_csv('bandwidth.csv', names=name)
    if not os.path.exists('./nodeinfo'):
        os.makedirs('./nodeinfo')
    i = 0
    while True :
        port = startPort + i
        cmdstr = '(./node -nport %d -bw %f -pscTimeAvg %d -pscTimeSig %d -nodeMode 1 -capa %f -machineType %d -groupIndex %d) 2>&1 | tee ./nodeinfo/node_%d.log &' % \
        (port, bw_data.at[i%50, 'bw'], timeAvg, timeAvg*6, capa, mType, gType, port)
        print(cmdstr)
        subprocess.call(cmdstr, shell=True)
        intervel = nextTime(nodeArrivalRate)
        time.sleep(intervel*60)
        i = i + 1


def stopNodes():
    cmdstr = "ps -ef|grep ./node|awk '{print $2}' | xargs kill "
    subprocess.call(cmdstr, shell=True)

def startClient(rate, trainNum, evalNum):
    time.sleep(1)
    cmdstr = '(./client -rate %d -ptrainNum %d -evalNum %d) 2>&1 | tee client.log' % (rate, trainNum, evalNum)
    subprocess.call(cmdstr, shell=True)

def SingleTest(nodeNum, dptIndex, timeAvg, rate, trainNum, evalNum, evalStr):
    pm = Process(target=startMaster, args=(dptIndex,))
    pn = Process(target=startNodes, args=(timeAvg,nodeNum, ))
    pc = Process(target=startClient, args=(rate, trainNum, evalNum))

    pm.start()
    pn.start()
    time.sleep(60*timeAvg) #half of timeAvg
    pc.start()
    pc.join()
    
    # clean up processes
    time.sleep(2)
    pm.terminate()
    pn.terminate()
    stopMaster()
    stopNodes()

    #save results
    time.sleep(2)
    testFileName = 'Result-%d-%d-%d-%s' % (dptIndex, timeAvg, rate, evalStr)
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

def SingleTestOnce(dptIndex, timeAvg, rate, trainNum, evalNum, nodeArrivalRate,evalStr):
    pm = Process(target=startMaster, args=(dptIndex,))
    pn = Process(target=startNodesOnce, args=(timeAvg,nodeArrivalRate))
    pc = Process(target=startClient, args=(rate, trainNum, evalNum))

    pn.daemon = True

    pm.start()
    pn.start()
    time.sleep(30*timeAvg) #half of timeAvg
    pc.start()
    pc.join()

    # clean up processes
    time.sleep(2)
    pm.terminate()
    pn.terminate()
    stopMaster()
    stopNodes()
 
    #save results
    time.sleep(2)

    testFileName = 'Result-%d-%d-%d-%d-%s' % (dptIndex, timeAvg, rate, nodeArrivalRate, evalStr)
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
    print('Evaluation is done. Dispatcher %d, timeAvg %d, taskRate %d nodeRate %d' % (dptIndex, timeAvg, rate, nodeArrivalRate))

def SingleTestHetero(dptIndex, taskRate, timeAvg1, nodeRate1, capa1, timeAvg2, nodeRate2, capa2, trainNum, evalNum):
    if not os.path.exists('./nodeinfo'):
        os.makedirs('./nodeinfo')
    pm = Process(target=startMaster, args=(dptIndex,))
    pn = Process(target=startNodesOnceWithCapa, args=(timeAvg1,nodeRate1, capa1, 0, 0, 50100))
    pn2 = Process(target=startNodesOnceWithCapa, args=(timeAvg2, nodeRate2, capa2, 1, 1, 55100))
    pc = Process(target=startClient, args=(taskRate, trainNum, evalNum))

    pn.daemon = True
    pn2.daemon = True

    pm.start()
    pn.start()
    pn2.start()
    time.sleep(30*timeAvg1) #half of timeAvg
    pc.start()
    pc.join()

    # clean up processes
    time.sleep(2)
    pm.terminate()
    pn.terminate()
    pn2.terminate()
    stopMaster()
    stopNodes()

    #save results
    time.sleep(2)
    testFileName = 'Result-H-%d-%d-%d-%d-%f-%d-%d-%f' % (dptIndex, taskRate, timeAvg1, nodeRate1, capa1, timeAvg2, nodeRate2, capa2)
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
    print('Evaluation is done. Dispatcher %d, taskRate %d, Node1: timeAvg %d, nodeRate %d capa:%f, Node2: timeAvg %d,nodeRate %d capa:%f' %\
         (dptIndex, taskRate, timeAvg1, nodeRate1, capa1, timeAvg2, nodeRate2, capa2))

if __name__ == "__main__":
    '''
    presenceTimeAvgList = [10, 15, 20, 25, 0]
    jobArrivalRateList = [90, 120, 150, 180, 210]
    disIndexList = [1, 2, 3, 4, 5]
    '''
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
    #SingleTest(5, 15, 90, 5000, 20000)
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
    '''
    rate = 120
    for timeAvg in presenceTimeAvgList:
        for index in disIndexList:
            SingleTest(index, timeAvg, rate, 5000, 20000)
            time.sleep(2)
    '''
    '''
    newDisIndexList = [6,7]
    timeAvg = 15
    for rate in jobArrivalRateList:
        for index in newDisIndexList:
            SingleTest(index, timeAvg, rate, 5000, 20000)
            time.sleep(2)

    rate = 120
    for timeAvg in presenceTimeAvgList:
        for index in newDisIndexList:
            SingleTest(index, timeAvg, rate, 5000, 20000)
            time.sleep(2)
    '''
    '''
    SingleTest(5, 15, 90, 5000, 20000)
    '''

    '''
    #改变 arrival rate， 30 ,60
        jobArrivalRateList = [30,60]
        disIndexList = [1, 6, 7, 5]
        timeAvg = 15
        for rate in jobArrivalRateList:
            for index in disIndexList:
                SingleTest(index, timeAvg, rate, 5000, 20000)
                time.sleep(2)
    '''

    #SingleTestOnce(1, 15, 120, 5000, 20000, 2)
    '''
    disIndexList = [6, 7, 5]
    nodeArrivalRateList = [2]
    timeAvg = 15
    jobRate = 120
    for nodeRate in nodeArrivalRateList:
        for index in disIndexList:
            SingleTestOnce(index, timeAvg, jobRate, 5000, 20000, nodeRate)
            time.sleep(2)

    disIndexList = [1, 6, 7, 5]
    nodeArrivalRateList = [4, 6, 8, 10]
    timeAvg = 15
    jobRate = 120
    for nodeRate in nodeArrivalRateList:
        for index in disIndexList:
            SingleTestOnce(index, timeAvg, jobRate, 5000, 20000, nodeRate)
            time.sleep(2)
    '''
    '''
    disIndexList = [8]
    nodeArrivalRateList = [2, 4, 6, 8, 10]
    timeAvg = 15
    jobRate = 120
    for nodeRate in nodeArrivalRateList:
        for index in disIndexList:
            SingleTestOnce(index, timeAvg, jobRate, 5000, 20000, nodeRate)
            time.sleep(2)
    '''

    #SingleTestHetero(dptIndex, taskRate, timeAvg1, nodeRate1, capa1, timeAvg2, nodeRate2, capa2, trainNum, evalNum)

    #SingleTestHetero(1, 120, 15, 5, 1.0, 5, 1, 1.5, 100, 100)

    #Heterogeneity
    '''
    node2RateList = [1]
    disIndexList = [6, 7, 5, 8]
    for node2Rate in node2RateList:
        for index in disIndexList:
            SingleTestHetero(index, 120, 15, 5, 1.0, 5, node2Rate, 1.5, 2500, 12500)

    node2RateList = [2, 3, 4, 5]
    disIndexList = [1, 6, 7, 5, 8]
    for node2Rate in node2RateList:
        for index in disIndexList:
            SingleTestHetero(index, 120, 15, 5, 1.0, 5, node2Rate, 1.5, 2500, 12500)
    '''
    '''
    #SingleTestHetero(index, 120, 15, 5, 1.0, 5, node2Rate, 1.5, 2500, 12500)
    SingleTestHetero(1, 120, 15, 5, 1.0, 5, 3, 1.5, 2500, 12500)
    SingleTestHetero(1, 120, 15, 5, 1.0, 5, 5, 1.5, 2500, 12500)
    SingleTestHetero(6, 120, 15, 5, 1.0, 5, 2, 1.5, 2500, 12500)
    SingleTestHetero(7, 120, 15, 5, 1.0, 5, 2, 1.5, 2500, 12500)
    SingleTestHetero(7, 120, 15, 5, 1.0, 5, 3, 1.5, 2500, 12500)
    SingleTestHetero(8, 120, 15, 5, 1.0, 5, 3, 1.5, 2500, 12500)
    SingleTestHetero(8, 120, 15, 5, 1.0, 5, 4, 1.5, 2500, 12500)
    SingleTestHetero(8, 120, 15, 5, 1.0, 5, 5, 1.5, 2500, 12500)
    '''

    '''
    presenceTimeAvgList = [10, 15, 20, 25, 30, 60, 0]
    #jobArrivalRateList = [90, 120, 150, 180, 210]
    disIndexList = [1, 5, 6, 7, 8]
    jobRate = 120

    for timeAvg in presenceTimeAvgList:
        for index in disIndexList:
            SingleTest(index, timeAvg, jobRate, 2500, 12500)
            time.sleep(2)
    '''
    '''
    node2RateList = [0]
    disIndexList = [6, 7, 5]
    for node2Rate in node2RateList:
        for index in disIndexList:
            SingleTestHetero(index, 120, 15, 5, 1.0, 5, node2Rate, 1.5, 2500, 12500)
            time.sleep(5)
    '''

    '''
    node2RateList = [1, 2, 3, 4, 5]
    disIndexList = [1, 6, 7, 5]
    for node2Rate in node2RateList:
        for index in disIndexList:
            SingleTestHetero(index, 120, 15, 5, 1.0, 5, node2Rate, 1.5, 2500, 12500)
            time.sleep(5)
    '''

    #SingleTest(nodeNum, dptIndex, timeAvg, rate, trainNum, evalNum, evalStr)

    '''
    presenceTimeAvgList = [10, 15, 20, 25, 30, 0]
    jobRate = 120
    disIndexList = [1, 6, 5, 7]
    for timeAvg in presenceTimeAvgList:
        for index in disIndexList:
            SingleTest(75, index, timeAvg, jobRate, 2500, 12500, 'presencetime')
            time.sleep(2)        

    presenceTimeAvgList = [25]
    jobRate = 120
    disIndexList = [7]
    for timeAvg in presenceTimeAvgList:
        for index in disIndexList:
            SingleTest(75, index, timeAvg, jobRate, 2500, 12500, 'presencetime')
            time.sleep(2)     

    presenceTimeAvgList = [30]
    jobRate = 120
    disIndexList = [1, 6, 5, 7]
    for timeAvg in presenceTimeAvgList:
        for index in disIndexList:
            SingleTest(75, index, timeAvg, jobRate, 2500, 12500, 'presencetime')
            time.sleep(2)   
    '''  

    '''
    presenceTimeAvgList = [20]
    jobRate = 90
    disIndexList = [1, 6, 5, 7]
    for timeAvg in presenceTimeAvgList:
        for index in disIndexList:
            SingleTest(75, index, timeAvg, jobRate, 2500, 12500, 'presencetime')
            time.sleep(2)          

    presenceTimeAvgList = [25]
    jobRate = 90
    disIndexList = [1, 6]
    for timeAvg in presenceTimeAvgList:
        for index in disIndexList:
            SingleTest(75, index, timeAvg, jobRate, 2500, 12500, 'presencetime')
            time.sleep(2)    
    '''
    '''
    jobArrivalRateList = [60, 120]
    disIndexList = [1,5,6,7,8]
    timeAvg = 20
    nodeRate = 5
    for jobRate in jobArrivalRateList:
        for index in disIndexList:
            SingleTestOnce(index, timeAvg, jobRate, 3000, 15000, nodeRate, 'taskRate')
            time.sleep(10)     
    '''

    '''
    presenceTimeAvgList = [15, 20, 25, 30, 0]
    jobRate = 180
    disIndexList = [1,5,6,7]
    for timeAvg in presenceTimeAvgList:
        for index in disIndexList:
            SingleTest(100, index, timeAvg, jobRate, 3000, 15000, 'presencetime')
            time.sleep(2)   
    '''

    '''
    jobArrivalRateList = [60, 120, 150, 180, 210, 240, 270, 300, 330, 360]
    disIndexList = [1,5,6,7]
    timeAvg = 20
    nodeRate = 5
    for jobRate in jobArrivalRateList:
        for index in disIndexList:
            SingleTestOnce(index, timeAvg, jobRate, 3000, 15000, nodeRate, 'taskRate')
            time.sleep(10)  
    '''
    
    '''
    jobArrivalRateList = [60]
    disIndexList = [6]
    timeAvg = 20
    nodeRate = 5
    for jobRate in jobArrivalRateList:
        for index in disIndexList:
            SingleTestOnce(index, timeAvg, jobRate, 3000, 15000, nodeRate, 'taskRate')
            time.sleep(10)  
    '''

    presenceTimeAvgList = [15, 20, 25, 30, 35, 40, 0]
    jobRate = 180
    disIndexList = [1,5,6,7]
    for timeAvg in presenceTimeAvgList:
        for index in disIndexList:
            SingleTest(100, index, timeAvg, jobRate, 3000, 15000, 'presencetime')
            time.sleep(2)  