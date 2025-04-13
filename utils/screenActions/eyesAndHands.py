import pyautogui
import os

def define_initial_region():
    screenX, screenY = pyautogui.size()
    return (0, 0, screenX, screenY)

def find_image(image_path, region=None):
    try:
        # If no region is provided, use the full screen
        if region is None:
            region = define_initial_region()
        
        location = pyautogui.locateOnScreen(image_path, confidence=0.8, region=region)
        if location is not None:
            x, y = pyautogui.center(location)
            pyautogui.moveTo(x, y)
            return location
    except Exception as e:
        print(f"Error in find_image: {e}")
        return None

def adjust_search_region(location, margin=50):
    left, top, width, height = location
    search_region = (
        max(0, left - margin),
        max(0, top - margin),
        width + (2 * margin),
        height + (2 * margin)
    )
    return search_region

if __name__ == "__main__":
    image_filepath = os.path.join(os.path.dirname(__file__), "imagesToFind", "test.png")
    search_region = None
    for i in range(10):
        location = find_image(image_filepath, region=search_region)
        if location is not None:
            search_region = adjust_search_region(location)
            print(search_region)
        else:
            print("No location found")
            location = find_image(image_filepath, region=define_initial_region())
