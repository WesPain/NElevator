package main

import "./elevio"
import "fmt"

func main(){

    numFloors := 4
    var presentFloor = -1
    var reqFloor int =-1
    var first bool =true

    elevio.Init("localhost:15657", numFloors)
    
    var d elevio.MotorDirection = elevio.MD_Up  //Elevator direction
    //elevio.SetMotorDirection(d)
    
    drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
    drv_stop    := make(chan bool)
    doorTimeout :=make(chan bool)
    doorTimerReset := make (chan bool)    
    
    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
    go elevio.PollStopButton(drv_stop)
    go elevio.DoorTimer(doorTimeout,doorTimerReset)
    
    
    for {
        select {
        case readButton := <- drv_buttons: //Button
            fmt.Printf("%+v\n", readButton)
            elevio.SetButtonLamp(readButton.Button, readButton.Floor, true)
            reqFloor = readButton.Floor
            //fmt.Printf("%+v\n", reqFloor)
            if readButton.Floor < presentFloor {
            	//fmt.Printf("test")
            	d = elevio.MD_Down
            } else if readButton.Floor > presentFloor {
            	d = elevio.MD_Up
            	
            } 
            elevio.SetMotorDirection(d)




            
        case sensorFloor := <- drv_floors: //Sensor
            fmt.Printf("%+v\n", sensorFloor)
            presentFloor = sensorFloor
            if first == true {

            	if (sensorFloor == 1 || sensorFloor==2 || sensorFloor ==3 || sensorFloor ==4){
            		d = elevio.MD_Stop
            		first = false
            	}
            	
            }
            if reqFloor== presentFloor{
            	d = elevio.MD_Stop
            	elvio.SetDoorOpenLamp(true)
            	doorTimerReset <- true
            }
            elevio.SetMotorDirection(d)
            //reqFloor= sensorFloor
            /*if sensorFloor == numFloors-1 {
                d = elevio.MD_Down
            } else if sensorFloor == 0 {
                d = elevio.MD_Up
            }
            elevio.SetMotorDirection(d)*/
            /*if first == true {
            	if (sensorFloor == 1 || sensorFloor==2 || sensorFloor ==3 || sensorFloor ==4){
            		d = elevio.MD_Stop
            	}
            	first = false
            }
            elevio.SetMotorDirection(d)*//*
            if reqFloor < sensorFloor {
            	//fmt.Printf("test")
            	d = elevio.MD_Down
            } else if reqFloor > sensorFloor {
            	d = elevio.MD_Up
            } else if sensorFloor == reqFloor{
            	d = elevio.MD_Stop
            }
            elevio.SetMotorDirection(d)*/
            
            
            
        case a := <- drv_obstr:
            fmt.Printf("%+v\n", a)
            if a {
                elevio.SetMotorDirection(elevio.MD_Stop)
            } else {
                elevio.SetMotorDirection(d)
            }
            
        case a := <- drv_stop:
            fmt.Printf("%+v\n", a)
            for f := 0; f < numFloors; f++ {
                for b := elevio.ButtonType(0); b < 3; b++ {
                    elevio.SetButtonLamp(b, f, false)
                }
            }
        case a := <- doorTimeout:
        	fmt.Printf("%+v\n", a)
			fmt.Printf("LOL")
			elvio.SetDoorOpenLamp(false)
			
        }
    }    
}
