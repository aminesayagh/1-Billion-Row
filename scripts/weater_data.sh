#!/bin/bash

# Filename for the output CSV file
OUTPUT_FILE="weather_data.csv"

# Required free storage in bytes (16GB)
REQUIRED_FREE_STORAGE=$((16 * 1024 * 1024 * 1024))

# Minimum required RAM in bytes (4GB)
REQUIRED_FREE_RAM=$((4 * 1024 * 1024 * 1024))

# Maximum allowed CPU load average (over 1 minute)
MAX_CPU_LOAD=2.0

# Array of station names
STATIONS=(
    Algeria Angola Benin Botswana "Burkina Faso" Burundi Cameroon "Cape Verde"
    "Central African Republic" Chad Comoros "Democratic Republic of the Congo"
    "Republic of the Congo" Djibouti Egypt "Equatorial Guinea" Eritrea Eswatini
    Ethiopia Gabon "The Gambia" Ghana Guinea "Guinea-Bissau" "Ivory Coast" Kenya
    Lesotho Liberia Libya Mauritania Morocco Namibia Niger Nigeria Rwanda
    "São Tomé and Príncipe" Senegal Seychelles "Sierra Leone" Somalia
    "South Africa" "South Sudan" Sudan Tanzania Togo Tunisia Uganda
    "Western Sahara" Zambia Zimbabwe "Antigua and Barbuda" Aruba Bahamas Barbados
    Belize Bermuda Canada "Cayman Islands" "Costa Rica" Cuba Dominica
    "Dominican Republic" "El Salvador" Greenland Grenada Guatemala Haiti
    Honduras Jamaica Mexico Montserrat "Netherlands Antilles" Nicaragua Panama
    "Saint Kitts and Nevis" "Saint Lucia" "Saint-Pierre and Miquelon"
    "Saint Vincent and the Grenadines" "Trinidad and Tobago" "Turks and Caicos Islands"
    "United States" "British Virgin Islands" "United States Virgin Islands" "Puerto Rico"
    Argentina Bolivia Brazil Chile Colombia Ecuador Guyana Paraguay Peru Suriname
    Uruguay Venezuela Afghanistan Azerbaijan Bahrain Bangladesh Bhutan Brunei Burma
    Cambodia China India Indonesia Iran Iraq Israel Japan Jordan Kazakhstan "North Korea"
    "South Korea" Kuwait Kyrgyzstan Laos Lebanon Malaysia Maldives Mongolia Nepal Oman
    Pakistan "The Philippines" Qatar Russia "Saudi Arabia" Singapore "Sri Lanka" Syria
    Taiwan Tajikistan Thailand Turkey Turkmenistan "United Arab Emirates" Uzbekistan
    Vietnam Yemen Albania Andorra Armenia Austria Azerbaijan Belarus Belgium
    "Bosnia and Herzegovina" Bulgaria Croatia Cyprus "Czech Republic" Denmark Estonia
    Finland France "Georgia (country)" Germany Greece Hungary Iceland Ireland Italy
    Kazakhstan Kosovo Latvia Liechtenstein Lithuania Luxembourg Macedonia Malta Moldova
    Montenegro Netherlands Norway Poland Portugal Romania Russia "San Marino" Serbia
    Slovakia Slovenia Spain Sweden Switzerland Turkey Ukraine "United Kingdom" Australia
    "East Timor" Fiji Indonesia Kiribati "Papua New Guinea" "Marshall Islands"
    "Federated States of Micronesia" Nauru "New Zealand" Palau Samoa "Solomon Islands"
    Tonga Tuvalu Vanuatu "Cook Islands" Niue "American Samoa" Guam "Northern Mariana Islands"
    "French Polynesia" "New Caledonia" "Wallis and Futuna" "Tuamotu" "Fiji Islands"
    "Rotuma Island" "Vanua Levu" "Viti Levu" "Kadavu Island" "Lau Island" "Lomaiviti Islands"
    "Mamanuca Islands" "Yasawa Islands" "Bau Island" "Beqa Island" "Denarau Island"
    "Kadavu Island" "Malolo Island" "Mana Island" "Matamanoa Island" "Monuriki Island"
    "Nananu-i-Ra Island" "Nukulau Island" "Qamea Island" "Tavarua Island" "Vatulele Island"
    "Viti Levu" "Yasawa Island" "Vatulele Island" "Viti Levu" "Yasawa Island" "Vatulele Island"
    "Viti Levu" "Yasawa Island" "Vatulele Island" "Viti Levu" "Yasawa Island" "Vatulele Island"
)

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check free storage
check_free_storage() {
    local available_storage
    available_storage=$(df --output=avail -k . | tail -n 1)
    available_storage=$((available_storage * 1024)) # Convert to bytes

    if (( available_storage < REQUIRED_FREE_STORAGE )); then
        echo "Error: Not enough free storage. At least 16GB of free storage is required."
        exit 1
    fi
}

# Function to check free RAM
check_free_ram() {
    local free_ram
    free_ram=$(free -b | awk '/Mem:/ {print $7}')

    if (( free_ram < REQUIRED_FREE_RAM )); then
        echo "Error: Not enough free RAM. At least 4GB of free RAM is required."
        exit 1
    fi
}

# Function to check CPU load
check_cpu_load() {
    local cpu_load
    cpu_load=$(uptime | awk -F'load average:' '{ print $2 }' | cut -d, -f1)

    if (( $(echo "$cpu_load > $MAX_CPU_LOAD" | bc -l) )); then
        echo "Error: CPU load is too high. Load average over 1 minute should be below $MAX_CPU_LOAD."
        exit 1
    fi
}

# Function to check write permissions
check_write_permissions() {
    if [ ! -w . ]; then
        echo "Error: No write permissions in the current directory."
        exit 1
    fi
}

# Function to generate a random temperature within the range [-99.9, 99.9]
generate_temperature() {
    local mean=15 # Average global temperature
    local stddev=30 # Standard deviation to allow for a wide range of temperatures

    # Generate a normally distributed random number
    local temp=$(awk -v mean="$mean" -v stddev="$stddev" '
        function box_muller(mean, stddev) {
            u1 = rand()
            u2 = rand()
            z0 = sqrt(-2 * log(u1)) * cos(2 * pi * u2)
            return mean + z0 * stddev
        }
        BEGIN {
            srand()
            printf "%.1f", box_muller(mean, stddev)
        }
    ')
    echo "$temp"
}

# Function to generate weather data
generate_weather_data() {
    local station
    local temperature
    local station_index
    local batch_size=10000
    local rows=1000000000-30000
    local batch=()

    echo "Generating weather data..."

    for ((i=1; i<=rows; i++)); do
        station_index=$((RANDOM % ${#STATIONS[@]}))
        station=${STATIONS[$station_index]}
        # temperature=$(generate_temperature)
        temperature=$(awk -v min=-99.9 -v max=99.9 'BEGIN { srand(); printf "%.1f", min + rand() * (max - min) }')
        batch+=("$station;$temperature")

        # Write batch to file
        if (( i % batch_size == 0 )); then
            printf "%s\n" "${batch[@]}" >> "$OUTPUT_FILE"
            batch=()
            echo "Generated $i rows..."
        fi
    done

    # Write any remaining data
    if [ ${#batch[@]} -gt 0 ]; then
        printf "%s\n" "${batch[@]}" >> "$OUTPUT_FILE"
    fi
}

# Check if seq, awk, and head are installed
if ! command_exists seq; then
    echo "Error: seq command not found. Please install it."
    exit 1
fi

if ! command_exists awk; then
    echo "Error: awk command not found. Please install it."
    exit 1
fi

if ! command_exists head; then
    echo "Error: head command not found. Please install it."
    exit 1
fi

# Perform all checks
check_free_storage
check_free_ram
# check_cpu_load
check_write_permissions


# Measure the time taken to generate the data
start_time=$(date +%s)

# Call the function
generate_weather_data

end_time=$(date +%s)
elapsed_time=$(( end_time - start_time ))

echo "Weather data generation complete. File saved as $OUTPUT_FILE."
echo "Time taken: $elapsed_time seconds."

# Print the first 10 lines of the generated file
echo "First 10 lines of the generated file:"
head -n 10 "$OUTPUT_FILE"
