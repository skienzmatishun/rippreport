import re

# Read the XML file as a string
with open("ripp.xml", "r") as f:
    xml_string = f.read()

# Define a regular expression pattern to find "wp:" in tag names
pattern = re.compile(r'<wp:(\w+)')

# Replace "wp:" with an empty string in all matches
xml_string = pattern.sub(r'<\1', xml_string)

# Write the modified XML to a new file
with open("ripp_no_prefix.xml", "w") as f:
    f.write(xml_string)
