package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

const (
	outputFile    = "weather_data.csv"
	requiredFreeStorage = 16 * 1024 * 1024 * 1024 // 16 GB
	requiredFreeRAM     = 4 * 1024 * 1024 * 1024 // 4 GB
	maxCPULoad          = 2.0
	rows                = 1000000000
	batchSize           = 100000
)

var stations = []string{
	"Algeria", "Angola", "Benin", "Botswana", "Burkina Faso", "Burundi", "Cameroon", "Cape Verde",
	"Central African Republic", "Chad", "Comoros", "Democratic Republic of the Congo",
	"Republic of the Congo", "Djibouti", "Egypt", "Equatorial Guinea", "Eritrea", "Eswatini",
	"Ethiopia", "Gabon", "The Gambia", "Ghana", "Guinea", "Guinea-Bissau", "Ivory Coast", "Kenya",
	"Lesotho", "Liberia", "Libya", "Mauritania", "Morocco", "Namibia", "Niger", "Nigeria", "Rwanda",
	"São Tomé and Príncipe", "Senegal", "Seychelles", "Sierra Leone", "Somalia", "South Africa",
	"South Sudan", "Sudan", "Tanzania", "Togo", "Tunisia", "Uganda", "Western Sahara", "Zambia", "Zimbabwe",
	"Antigua and Barbuda", "Aruba", "Bahamas", "Barbados", "Belize", "Bermuda", "Canada", "Cayman Islands",
	"Costa Rica", "Cuba", "Dominica", "Dominican Republic", "El Salvador", "Greenland", "Grenada",
	"Guatemala", "Haiti", "Honduras", "Jamaica", "Mexico", "Montserrat", "Netherlands Antilles",
	"Nicaragua", "Panama", "Saint Kitts and Nevis", "Saint Lucia", "Saint-Pierre and Miquelon",
	"Saint Vincent and the Grenadines", "Trinidad and Tobago", "Turks and Caicos Islands", "United States",
	"British Virgin Islands", "United States Virgin Islands", "Puerto Rico", "Argentina", "Bolivia", "Brazil",
	"Chile", "Colombia", "Ecuador", "Guyana", "Paraguay", "Peru", "Suriname", "Uruguay", "Venezuela",
	"Afghanistan", "Azerbaijan", "Bahrain", "Bangladesh", "Bhutan", "Brunei", "Burma", "Cambodia", "China",
	"India", "Indonesia", "Iran", "Iraq", "Israel", "Japan", "Jordan", "Kazakhstan", "North Korea",
	"South Korea", "Kuwait", "Kyrgyzstan", "Laos", "Lebanon", "Malaysia", "Maldives", "Mongolia", "Nepal",
	"Oman", "Pakistan", "The Philippines", "Qatar", "Russia", "Saudi Arabia", "Singapore", "Sri Lanka",
	"Syria", "Taiwan", "Tajikistan", "Thailand", "Turkey", "Turkmenistan", "United Arab Emirates", "Uzbekistan",
	"Vietnam", "Yemen", "Albania", "Andorra", "Armenia", "Austria", "Azerbaijan", "Belarus", "Belgium",
	"Bosnia and Herzegovina", "Bulgaria", "Croatia", "Cyprus", "Czech Republic", "Denmark", "Estonia",
	"Finland", "France", "Georgia (country)", "Germany", "Greece", "Hungary", "Iceland", "Ireland", "Italy",
	"Kazakhstan", "Kosovo", "Latvia", "Liechtenstein", "Lithuania", "Luxembourg", "Macedonia", "Malta",
	"Moldova", "Montenegro", "Netherlands", "Norway", "Poland", "Portugal", "Romania", "Russia", "San Marino",
	"Serbia", "Slovakia", "Slovenia", "Spain", "Sweden", "Switzerland", "Turkey", "Ukraine", "United Kingdom",
	"Australia", "East Timor", "Fiji", "Indonesia", "Kiribati", "Papua New Guinea", "Marshall Islands",
	"Federated States of Micronesia", "Nauru", "New Zealand", "Palau", "Samoa", "Solomon Islands", "Tonga",
	"Tuvalu", "Vanuatu",
}

func generateTemperature() float64 {
	mean := 15.0
	stddev := 30.0
	u1 := rand.Float64() // Generate a random number between 0 and 1
	u2 := rand.Float64() // Generate a random number between 0 and 1
	z0 := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2.0*math.Pi*u2) // Box-Muller transform, generates a random number from a normal distribution, 
	return mean + z0*stddev // Return the temperature
}

func main() {
	startTime := time.Now()
	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer output.Close()

	fmt.Println("Generating weather data...")

	batch := make([]any, 0, batchSize)
	for i := 1; i < rows; i++ {
		station := stations[rand.Intn(len(stations))] // Randomly select a station
		temperature := generateTemperature() 		// Generate a random temperature
		batch = append(batch, fmt.Sprintf("%s;%.2f\n", station, temperature)) // Append the data to the batch

		// If the batch is full, write it to the file
		if i%batchSize == 0 { // Write the batch to the file
			output.WriteString(fmt.Sprintln(batch...))
			batch = batch[:0] // Clear the batch
			fmt.Printf("Generated %d rows...\n", i)
		}
	}

	// Write any remaining data
	if len(batch) > 0 {
		output.WriteString(fmt.Sprintln(batch...))
	}

	elapsedTime := time.Since(startTime) // Calculate the time taken
	fmt.Printf("Weather data generation complete. File saved as %s.\n", outputFile)
	fmt.Printf("Time taken: %f seconds.\n", elapsedTime.Seconds())
}
