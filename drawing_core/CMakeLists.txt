cmake_minimum_required(VERSION 3.20)
project("CelebratoryBullets")

# set flags
set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS ON)
set(CMAKE_CXX_FLAGS "-O3 -Wno-deprecated-declarations -Wno-unused-command-line-argument -Wno-cpp -lglfw -framework OpenGL")
set(CMAKE_PREFIX_PATH "/opt/homebrew")

# set(CMAKE_VERBOSE_MAKEFILE ON)

# set output dirs
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin)

# include project headers
include_directories(${CMAKE_SOURCE_DIR}/include)

# *import essentials from homebrew
include_directories(/opt/homebrew/include)
link_directories(/opt/homebrew/lib)

# *import from APPLE frameworks
include_directories(/Library/Developer/CommandLineTools/SDKs/MacOSX.sdk/System/Library/Frameworks)
FIND_LIBRARY(COCOA_LIBRARY Cocoa)
FIND_LIBRARY(OpenGL_LIBRARY OpenGL)
FIND_LIBRARY(IOKit_LIBRARY IOKit)
FIND_LIBRARY(CoreVideo_LIBRARY CoreVideo)
MARK_AS_ADVANCED(COCOA_LIBRARY OpenGL_LIBRARY)
SET(APPLE_LIBS ${COCOA_LIBRARY} ${IOKit_LIBRARY} ${OpenGL_LIBRARY} ${CoreVideo_LIBRARY})

# SET(APPLE_LIBS ${APPLE_LIBS} ${GLFW3_LIBRARY} ${ASSIMP_LIBRARY} ${FREETYPE_LIBRARIES})

# create and link all main libraries into LIBS
set(ALL_LIBS ${ALL_LIBS} ${APPLE_LIBS} "src/shader.cpp" "src/stb_image.cpp")
add_library(LIBS SHARED ${ALL_LIBS})
link_libraries(LIBS)

add_subdirectory(src)