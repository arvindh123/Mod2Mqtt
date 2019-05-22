#sleep 6
python3 /home/pi/timeupdate.py > /dev/null 2>&1 &
python3  /home/pi/codeV0.1.py > /dev/null 2>&1 &
#python3  /home/pi/TextFileSplitter.py > /dev/null 2>&1 &

#sleep 30
#python3 /home/pi/Bovone/textread.py > /dev/null 2>&1 &


/home/pi/MultiPost/main > /home/pi/MultiPost/main.log 2>&1 &

./main > /home/pi/MultiPost/main.log 2>&1 &

 