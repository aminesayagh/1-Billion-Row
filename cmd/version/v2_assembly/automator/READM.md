# Automaton for floating point number parsing

This automaton is designed to minimize the number of states and transitions on the parsing of floating point numbers.

## States and Transitions

### States

- **q0 (Start):** The initial state where the automaton begins processing.
- **q1 (Sign):** Handles the optional sign (`+` or `-`) of the number.
- **q2 (Integer Part):** Processes digits that form the integer part of the number.
- **q3 (Decimal Point):** Represents the transition to the fractional part.
- **q4 (Leading Zero in Fraction):** Processes the initial `0` in the fractional part.
- **q5 (Fractional Part):** Processes digits after the decimal point.
- **q6 (Accept):** The accepting state, indicating a valid floating-point number.

### Transitions

- **q0 -> q1:** On input `+` or `-`, the automaton transitions from `q0` to `q1`.
- **q0 -> q2:** On input `1..9`, the automaton transitions from `q0` to `q2`.
- **q1 -> q2:** On input `1..9`, the automaton transitions from `q1` to `q2`.
- **q2 -> q2:** On input `0..9`, the automaton remains in `q2` (loop for multiple digits).
- **q2 -> q3:** On input `.`, the automaton transitions from `q2` to `q3`.
- **q3 -> q4:** On input `0`, the automaton transitions from `q3` to `q4`.
- **q3 -> q5:** On input `1..9`, the automaton transitions from `q3` to `q5`.
- **q4 -> q4:** On input `0..9`, the automaton remains in `q4`.
- **q4 -> q5:** On input `1..9`, the automaton transitions from `q4` to `q5`.
- **q5 -> q4:** On input `0`, the automaton transitions from `q5` back to `q4`.
- **q5 -> q5:** On input `1..9`, the automaton remains in `q5`.
- **q5 -> q6:** On input [accepting condition], the automaton transitions to the accepting state `q6`.

### Diagram

![Automaton Diagram](automaton.png)
*This diagram illustrates the state transitions within the automaton for processing floating-point numbers.*

## Optimizations

This automaton has to be optimized by the following techniques:

