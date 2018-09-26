#!/bin/bash
# Runs the blast benchmark for increasing numbers of requests. 

# Location of the results
RESULTS="throughput.json"

RUNS=12
MIN_OPS=50
MAX_OPS=1000 
OPS_INCR=50

# Describe the time format
TIMEFORMAT="experiment completed in %2lR"

time {
  # Run the experiment for each blast $RUNS times
  for (( I=0; I<$RUNS; I+=1 )); do

      # Run a benchmark from min ops to max ops by ops incr 
      for (( J=$MIN_OPS; J<=$MAX_OPS; J=J+$OPS_INCR )); do

        UPTIME=2500ms

        # Run Server
        speedmap serve -u $UPTIME & 
        sleep 1
        sclient blast -r $J >> $RESULTS
        wait

      done
  done
}
