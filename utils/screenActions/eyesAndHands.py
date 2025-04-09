import pyautogui
import os

def find_image(image_path, region):
    # Search within the specified region
    location = pyautogui.locateOnScreen(image_path, confidence=0.5, region=region)
    if location is not None:
        print("Image found")
        x, y = pyautogui.center(location)
        pyautogui.moveTo(x, y)
        return location  # Return the location if found
    else:
        print("Not found, searching again...")

def define_region():
    screenX, screenY = pyautogui.size()

def smaller_region(location, margin=50):
    # Expand the region slightly around the last known location
    left, top, width, height = location
    screenX, screenY = pyautogui.size()
    new_left = max(0, left - margin)
    new_top = max(0, top - margin)
    new_width = min(screenX, left + width + margin) - new_left
    new_height = min(screenY, top + height + margin) - new_top
    return (new_left, new_top, new_width, new_height)

def main():
    image_path = os.path.join(os.path.dirname(__file__), "imagesToFind", "test.png")
    region = define_region()

    for i in range(50):
        try:
            location = find_image(image_path, region)
            if location:
                region = smaller_region(location)  # Update the region dynamically
            else:
                print("Continuing search in the same region...")  # Keep the current region
        except Exception as e:
            print(f"An error occurred: {e}. Continuing...")  # Handle unexpected errors and continue

if __name__ == "__main__":
    main()