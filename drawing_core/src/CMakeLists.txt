get_filename_component(FOLDER_NAME ${CMAKE_CURRENT_SOURCE_DIR} NAME)
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin)
link_libraries(assimp freetype)
include_directories(/opt/homebrew/include/freetype2)
add_executable(fireworkbullet fireworkbullet.cpp)