import xml.etree.ElementTree as ET

# Parse the WordPress export file
tree = ET.parse('ripp.xml')
root = tree.getroot()

# Open the output file for writing
with open('output.txt', 'w') as f:
    # Loop through each post in the export file
    for item in root.iter('item'):
        # Extract the post title
        title = item.find('title').text

        # Extract the divs with id at_zurlpreview and their contents
        at_zurlpreview_divs = item.findall(".//div[@id='at_zurlpreview']")
        for div in at_zurlpreview_divs:
            # Extract the h3 heading and img source URL
            heading = div.find('h3').text or ''
            img_src = div.find('img').get('src') or ''

            # Write the heading and img src to the output file
            f.write(f"{title}:\n")
            f.write(f"\t{heading}\n")
            f.write(f"\t{img_src}\n")
