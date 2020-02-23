#/bin/bash
rm ./nodeinfo/*
#start_port=50100

#for i in $(seq 50101 50150);
for i in $(seq 50101 50120);
do 
    (./node -nport $i) 2>&1 | tee ./nodeinfo/node_$i.log &
    #(./node -nport $i -groupIndex $i ) 2>&1 | tee ./nodeinfo/node_$i.log &
done
