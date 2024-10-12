# z3b Encoding Specification

**Author:** Jaakko Heusala  
**Date:** October 12, 2024  
**Version:** 0.1

---

## Table of Contents

- [Abstract](#abstract)
- [Introduction](#introduction)
- [Motivation](#motivation)
- [Terminology](#terminology)
- [Character Set](#character-set)
- [Encoding Algorithm](#encoding-algorithm)
- [Decoding Algorithm](#decoding-algorithm)
- [Implementation Details](#implementation-details)
- [Examples](#examples)
- [Security Considerations](#security-considerations)
- [References](#references)
- [Acknowledgments](#acknowledgments)

---

## Abstract

This document specifies the **z3b encoding**, a binary-to-text encoding scheme
designed to map all 256 possible byte values to a set of printable ASCII 
characters safe for use in JSON strings without additional escaping. The
encoding ensures that every byte value is uniquely representable, and the 
resulting encoded strings exclude problematic characters such as the double 
quote (`"`) and control characters.

---

## Introduction

Data encoding schemes are essential for transmitting binary data over media 
designed to handle textual data. Common encoding schemes like Base64 are widely 
used but may not meet specific requirements such as exclusion of certain 
characters or optimizing for JSON compatibility.

The **z3b encoding** addresses these concerns by:

- Mapping all 256 byte values to a carefully selected set of printable ASCII characters.
- Ensuring the encoded data contains no characters that require escaping in JSON strings.
- Providing a reversible and unambiguous encoding and decoding process.

---

## Motivation

When embedding binary data within JSON strings, certain characters can cause 
issues due to JSON's syntax requirements and escaping mechanisms. The double 
quote (`"`) and control characters, in particular, necessitate escaping, which
can complicate data handling and increase payload size.

The **z3b encoding** is motivated by the need for an encoding scheme that:

- Avoids characters requiring escaping in JSON.
- Provides efficient encoding and decoding operations.
- Ensures that all possible byte values can be represented.

---

## Terminology

- **Byte**: An 8-bit unsigned integer ranging from 0 to 255.
- **Character**: A printable ASCII character used in the encoded output.
- **Set**: A collection of byte-to-character mappings.
- **Separator**: A special character used to switch between sets during encoding and decoding.

---

## Character Set

The encoding uses a specific set of printable ASCII characters, excluding the
double quote (`"`) to avoid JSON escaping issues.

**Printable Set (93 characters):**

```
!#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^`abcdefghijklmnopqrstuvwxyz{|}~
```

- **Total Characters:** 93
- **Excluded Characters:** `"`, space (` `)

---

## Encoding Algorithm

### Overview

- **Sets:** The encoding uses three sets (Set 1, Set 2, Set 3) to map byte values to characters.
- **Unique Mappings:** Each byte value from 0 to 255 is uniquely assigned to a character in one of the sets.
- **Set Switching:** A separator character (`_`) is used to switch between sets during encoding.

### Steps

1. **Initialize Sets:**
   - **Set 1:** Byte values 0x00 to 0x55 (0 to 85).
   - **Set 2:** Byte values 0x56 to 0xAB (86 to 171).
   - **Set 3:** Byte values 0xAC to 0xFF (172 to 255).

2. **Create Mappings:**
   - Map each byte value in a set to a character from the printable set.
   - The same printable characters are reused across sets.

3. **Encoding Process:**
   - Start with **Set 1**.
   - For each byte in the input data:
     - Attempt to map the byte using the current set.
     - If the byte is not found in the current set:
       - Append the **separator** (`_`) to the output.
       - Switch to the next set (Set 2 or Set 3).
       - Retry mapping the byte.
       - Repeat this process until the byte is successfully mapped.
     - Append the mapped character to the output.

---

## Decoding Algorithm

### Overview

- The decoding process reverses the encoding, using the same sets and mappings.
- The separator character indicates when to switch sets during decoding.

### Steps

1. **Initialize:**
   - Start with **Set 1**.
   - Prepare an empty output byte slice.

2. **Decoding Process:**
   - Iterate over each character in the encoded input:
     - If the character is the **separator** (`_`):
       - Switch to the next set (Set 2 or Set 3).
       - Continue to the next character.
     - Else:
       - Use the current set to map the character back to its corresponding byte value.
       - Append the byte to the output byte slice.
     - If the character cannot be mapped in the current set:
       - Return an error indicating an invalid character.

---

## Implementation Details

### Separator Character

- **Character:** `_` (underscore)
- **Purpose:** Indicates a switch to the next set during encoding and decoding.
- **Usage:** Appended to the encoded string before switching sets.

### Invalid Characters

- The encoding excludes certain characters to maintain JSON compatibility.
- **Invalid Character:** `"` (double quote)
- Any occurrence of an invalid character in the encoded data results in a decoding error.

### Data Structures

- **byteToCharSet:** A 2D array mapping each byte value to a character for each set.
- **charToByteSet:** A 2D array mapping each character back to its byte value for each set.
- **Sets:** Arrays `binarySet1`, `binarySet2`, and `binarySet3` containing the byte values assigned to each set.

---

## Examples

### Encoding Example

**Input Data (Hex):**

```
0x48 0x65 0x6C 0x6C 0x6F (ASCII "Hello")
```

**Encoding Steps:**

1. **Start with Set 1.**
2. **Byte 0x48:**
   - Not in Set 1.
   - Append `_` to output, switch to Set 2.
   - 0x48 is in Set 2, maps to character `!`.
   - Append `!` to output.
3. **Byte 0x65:**
   - 0x65 is in Set 2, maps to character `#`.
   - Append `#` to output.
4. **Byte 0x6C:**
   - 0x6C is in Set 2, maps to character `$`.
   - Append `$` to output.
5. **Byte 0x6C:**
   - 0x6C is in Set 2, maps to character `$`.
   - Append `$` to output.
6. **Byte 0x6F:**
   - 0x6F is in Set 2, maps to character `%`.
   - Append `%` to output.

**Encoded Output:**

```
_!#$$%
```

### Decoding Example

**Encoded Data:**

```
_!#$$%
```

**Decoding Steps:**

1. **Start with Set 1.**
2. **Character `_`:**
   - Separator, switch to Set 2.
3. **Character `!`:**
   - Maps to byte 0x48 in Set 2.
   - Append 0x48 to output.
4. **Character `#`:**
   - Maps to byte 0x65 in Set 2.
   - Append 0x65 to output.
5. **Character `$`:**
   - Maps to byte 0x6C in Set 2.
   - Append 0x6C to output.
6. **Character `$`:**
   - Maps to byte 0x6C in Set 2.
   - Append 0x6C to output.
7. **Character `%`:**
   - Maps to byte 0x6F in Set 2.
   - Append 0x6F to output.

**Decoded Output (Hex):**

```
0x48 0x65 0x6C 0x6C 0x6F
```

---

## Security Considerations

- **Data Integrity:** The encoding provides a reversible mapping, ensuring data integrity during encoding and decoding.
- **Character Restrictions:** By excluding problematic characters, the encoding reduces the risk of injection attacks when data is embedded in JSON strings.
- **Error Handling:** Implementations should handle decoding errors gracefully, particularly when encountering invalid characters or sequences.

---

## References

- [RFC 4648](https://tools.ietf.org/html/rfc4648): The Base16, Base32, and Base64 Data Encodings.
- [JSON Data Interchange Format](https://www.json.org/json-en.html)

---

## Acknowledgments

This encoding scheme was developed by **Jaakko Heusala** to address specific requirements for embedding binary data within JSON strings without the need for additional character escaping.

---

## Appendix: Implementation Notes

- **Language:** The encoding can be implemented in any programming language that supports basic data structures and byte manipulation.
- **Go Implementation:** A reference implementation is available in Go, utilizing arrays for mapping and optimized for performance.
- **Performance Considerations:**
  - Use of `[]byte` instead of `string` can enhance performance by reducing memory allocations.
  - Preallocating buffers based on input size estimates can minimize dynamic memory operations.

---

## Appendix: Test Vectors

**Test Case 1: Empty Input**

- **Input:** `[]`
- **Encoded Output:** `""`
- **Decoded Output:** `[]`

**Test Case 2: All Byte Values**

- **Input:** `[0x00, 0x01, ..., 0xFE, 0xFF]` (256 bytes)
- **Encoded Output:** (A string consisting of mapped characters and separators)
- **Decoded Output:** `[0x00, 0x01, ..., 0xFE, 0xFF]`

---

## Appendix: Character Mappings

Due to the length constraints, the full mapping tables are not included here. Implementers should ensure that:

- Each set's byte-to-character mapping is correctly established.
- The separator character (`_`) is not included in the printable set or mapped to any byte value.
- The invalid character (`"`) is not included in the printable set to maintain JSON compatibility.

---

## Appendix: Go Code Snippet

```go
// Example function signatures in Go

// Encode encodes the input bytes into a z3b-encoded byte slice.
func Encode(data []byte) ([]byte, error)

// Decode decodes a z3b-encoded byte slice back into bytes.
func Decode(encoded []byte) ([]byte, error)
```

---

## Contact Information

For questions or feedback regarding this encoding scheme, please contact:

**Jaakko Heusala**  
Email: [jheusala@iki.fi](mailto:jheusala@iki.fi)

---

*This document is licensed under the FSL-1.1-MIT License.*
