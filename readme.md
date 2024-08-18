# One Billion Data Parsing Project

<!-- a menu index can help -->

## Table of Contents

- [Introduction](#introduction)
- [How Much is a Billion?](#how-much-is-a-billion)
- [The Problem Statement](#the-problem-statement)
- [Hardware Constraints](#hardware-constraints)
- [Software Constraints](#software-constraints)
- [Software Input](#software-input)
  - [File Data](#file-data)
  - [Environment Variables Configuration](#environment-variables-configuration)
- [Criteria for Success](#criteria-for-success)
- [The Solution](#the-solution)
- [The Data](#the-data)

## Introduction

As Jockie Stewart once said, "You don't have to be an engineer to be a racing driver, but you do have to have mechanical sympathy." This project focuses on the parsing of a CSV file containing one billion rows of data. It's not just about writing code that works; it's about designing software that respects the hardware and the necessity of scale.

## How Much is a Billion?

A billion is a massive number. To count from 1 to 1 billion, it would take you 31 years, 259 days, 1 hour, 46 minutes, and 40 seconds. If you were to stack a billion pennies, it would reach a height of 870 miles. Driving a billion miles would allow you to circumnavigate the Earth 40,000 times. Clearly, a billion is a significant figure.

This scale is often required by large companies like Google, which has 1.2 billion users. They must store and process vast amounts of data from these users. Dealing with this scale presents a considerable challenge, and this project serves as a small example of how to tackle it.

## The Problem Statement

The problem is straightforward: you have a CSV file with one billion rows, and your task is to parse it. The CSV file follows the format:

```
<station_name:string>;<temperature:float>
```

You need to parse the file into a list where each row represents a station with its minimum temperature, maximum temperature, and average temperature. The list should be sorted by the station name and have the following format:

```
<station_name:string>;<min_temperature:float>;<max_temperature:float>;<medium_temperature:float>;<count_station:int>
```

## Hardware Constraints

Here are the hardware constraints for this project:

| Hardware Component | Constraint |
|--------------------|------------|
| CPU                | 8          |
| RAM                | 32GB       |
| Storage            | 1TB SSD    |
| Architecture       | x86_64     |
| Goroutines         | 1          |
| GoMaxProcs         | 8          |

## Software Constraints

Here are the software constraints for this project:

| Software Component | Constraint    |
|--------------------|---------------|
| Go Version         | 1.22.5        |
| Go Arch            | amd64         |
| File System        | ext4          |
| Operating System   | Linux         |
| Number of Files    | 1             |
| File Size          | 16GeGa        |
| File Format        | CSV           |
| File Columns       | 2             |
| File Rows          | 1,000,010,000 |
| File Encoding      | UTF-8         |
| File Delimiter     | `;`           |
| File Compression   | None          |
| File Line Ending   | `\n`          |

## Software Input

### File Data

### Environment Variables Configuration

## How to Run the Project

### Prerequisites

### Environment Variables

### Running the Project

## Criteria for Success

## The Solution

## The Data
