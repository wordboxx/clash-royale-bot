import pyautogui
import os

def find_image(image_path, region):
    try:
        # Search within the specified region
        location = pyautogui.locateOnScreen(image_path, confidence=0.8, region=region)
        if location is not None:
            print("Image found")
            x, y = pyautogui.center(location)
            pyautogui.moveTo(x, y)
    except Exception as e:
        print(f"Error in find_image: {e}")
        find_image(image_path, get_full_screen_region())


def get_full_screen_region():
    screenX, screenY = pyautogui.size()
    return (0, 0, screenX, screenY)

def adjust_search_region(location, screenX, screenY, margin=50):
    # Expand the region slightly around the last known location
    # (Speeds up the search by reducing the area)
    left, top, width, height = location
    new_left = max(0, left - margin)
    new_top = max(0, top - margin)
    new_width = min(screenX, left + width + margin) - new_left
    new_height = min(screenY, top + height + margin) - new_top
    return (new_left, new_top, new_width, new_height)

if __name__ == "__main__":
    image_path = os.path.join(os.path.dirname(__file__), "imagesToFind", "test.png")
    for i in range(10):
        find_image(image_path, get_full_screen_region())