import pyautogui
import os

def find_image(image_path, region):
    # Search within the specified region
    location = pyautogui.locateOnScreen(image_path, confidence=0.5, region=region)
    if location is not None:
        print("Image found")
        x, y = pyautogui.center(location)
        pyautogui.moveTo(x, y)
        return location
    else:
        print("Not found, searching again...")

def define_region():
    screenX, screenY = pyautogui.size()

def adjust_search_region(location, screenX, screenY, margin=50):
    # Expand the region slightly around the last known location
    # (Speeds up the search by reducing the area)
    left, top, width, height = location
    new_left = max(0, left - margin)
    new_top = max(0, top - margin)
    new_width = min(screenX, left + width + margin) - new_left
    new_height = min(screenY, top + height + margin) - new_top
    return (new_left, new_top, new_width, new_height)

def trackImage():
    image_path = os.path.join(os.path.dirname(__file__), "imagesToFind", "test.png")
    # TODO: region seems redundant, fix this in larger codebase
    region = define_region()
    screenX, screenY = pyautogui.size()  # Get screen size once

    for i in range(50):
        try:
            location = find_image(image_path, region)
            if location:
                region = adjust_search_region(location, screenX, screenY)  # Pass screen dimensions
            else:
                print("Continuing search in the same region...")  # Keep the current region
        except Exception as e:
            # Handle unexpected errors and continue if image not found
            # This will prevent the program from crashing
            print(f"An error occurred: {e}. Continuing...")