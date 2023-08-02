// Add support for a kelvin flag

package main

import  (
	"fmt"
	"flag"
	"strconv"
)

type KelvinTemp float64

type kelvinFlag struct { // define a concrete type that satisfies value interface of flag
	Kelvin KelvinTemp
}

func (k *kelvinFlag) String() string {
	return fmt.Sprintf("%f Kelvin", k.Kelvin)
}

func (k *kelvinFlag) Set(s string) error { // no conversions or units, value is straightaway expected in kelvin
	temp,err := strconv.ParseFloat(s, 64)
	k.Kelvin = KelvinTemp(temp)
	return err
}


func main() {
	kelvin := KelvinFlag("kTemp", KelvinTemp(123), "Used to indicate the temperature in Kelvin")
	flag.Parse()
	fmt.Println(*kelvin)

}

func KelvinFlag(flagName string, val KelvinTemp, helpStr string) *KelvinTemp {
	k := &kelvinFlag{val}
	flag.CommandLine.Var(k, flagName, helpStr)
	return &k.Kelvin
}



