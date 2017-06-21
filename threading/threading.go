package main

/*
#define _GNU_SOURCE
#include <sched.h>
#include <pthread.h>
#include <signal.h>
#include <unistd.h>
#include <stdio.h>

unsigned int
pthreadSelf()
{
	pthread_t tid;
	tid = pthread_self();
	return tid;
}

int
setAffinity(int cpu)
{
	cpu_set_t cpuset;
	CPU_ZERO(&cpuset);
	CPU_SET(cpu , &cpuset );
	return sched_setaffinity(0, sizeof(cpuset), &cpuset);
}

int
checkAffinity(int cpu)
{
	int ret;
	cpu_set_t cpuset, desired_cpuset;


	CPU_ZERO(&desired_cpuset);
	CPU_SET(cpu , &desired_cpuset );

	CPU_ZERO(&cpuset);
	CPU_SET(cpu , &cpuset ); //set CPU 2 on cpuset
	ret = sched_getaffinity(0, sizeof(cpuset), &cpuset);
	if (ret != 0) {
		return -1;
	}

	if (CPU_EQUAL(&cpuset, &desired_cpuset) == 0) {
		return -1;
	}
	return 0;
}

*/
import "C"

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func myThread(name string, wg *sync.WaitGroup, cpu int, setAffinity bool, lockGoThread bool, expectLock bool) {
	defer wg.Done()

	if lockGoThread {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
	}

	time.Sleep(1 * time.Second)

	if setAffinity {
		C.setAffinity(C.int(cpu))
	}

	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)

		tid := C.pthreadSelf()
		//fmt.Printf("[%v] Thread ID [%s] := %v \n", i, name, tid)
		if C.checkAffinity(C.int(cpu)) != 0 {
			if expectLock {
				fmt.Printf("FATAL: Affinity failure [%v] Thread ID [%s] := %v \n", i, name, tid)
			} else {
				fmt.Printf("Affinity failure [%v] Thread ID [%s] := %v \n", i, name, tid)
			}
		}
	}

}

func main() {
	var wg sync.WaitGroup
	tid := C.pthreadSelf()
	fmt.Println("My thread id :=", tid)
	wg.Add(4)
	go myThread("first", &wg, 1, true, true, true)
	go myThread("second", &wg, 2, false, true, false)
	go myThread("third", &wg, 3, true, false, false)
	go myThread("fourth", &wg, 4, false, false, false)
	wg.Wait()
}
