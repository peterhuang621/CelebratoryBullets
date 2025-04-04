#include <GLFW/glfw3.h>
#include <atomic>
#include <iostream>
#include <map>
#include <memory>
#include <string>

#include <glm/glm.hpp>
#include <glm/gtc/matrix_transform.hpp>
#include <glm/gtc/type_ptr.hpp>

#include <freetype2/ft2build.h>
#include FT_FREETYPE_H

#include "camera.h"
#include "shader.h"
#include "stb_image.h"
using namespace std;
const unsigned int SCR_WIDTH = 800;
const unsigned int SCR_HEIGHT = 600;

Camera camera(glm::vec3(0.0f, 0.0f, 5.0f));
float lastX = (float)SCR_WIDTH / 2.0;
float lastY = (float)SCR_HEIGHT / 2.0;
bool firstMouse = true;

float deltaTime = 0.0f;
float lastFrame = 0.0f;

struct Character
{
    unsigned int TextureID; // ID handle of the glyph texture
    glm::ivec2 Size;        // Size of glyph
    glm::ivec2 Bearing;     // Offset from baseline to left/top of glyph
    unsigned int Advance;   // Horizontal offset to advance to next glyph
};

struct Bullet
{
    float DurationSecs, Size;
    float Color[4], Position[3];
    string s;
    Bullet(float a, float b, float c[4], float d[3], const string &ss) : DurationSecs(a), Size(b), s(ss)
    {
        for (int j = 0; j < 4; j++)
            Color[j] = c[j];
        for (int j = 0; j < 3; j++)
            Position[j] = d[j];
    }
    Bullet() : DurationSecs(1.0), Size(1.0), s("NULL")
    {
        for (int j = 0; j < 4; j++)
            Color[j] = 0.5;
        for (int j = 0; j < 3; j++)
            Position[j] = 50.0;
    }
};

map<GLchar, Character> Characters;
unsigned int VAO, VBO;

void processInput(GLFWwindow *window)
{
    if (glfwGetKey(window, GLFW_KEY_ESCAPE) == GLFW_PRESS)
        glfwSetWindowShouldClose(window, true);
    if (glfwGetKey(window, GLFW_KEY_W) == GLFW_PRESS)
        camera.ProcessKeyboard(FORWARD, deltaTime);
    if (glfwGetKey(window, GLFW_KEY_S) == GLFW_PRESS)
        camera.ProcessKeyboard(BACKWARD, deltaTime);
    if (glfwGetKey(window, GLFW_KEY_A) == GLFW_PRESS)
        camera.ProcessKeyboard(LEFT, deltaTime);
    if (glfwGetKey(window, GLFW_KEY_D) == GLFW_PRESS)
        camera.ProcessKeyboard(RIGHT, deltaTime);
}

void framebuffer_size_callback(GLFWwindow *window, int width, int height)
{
    glViewport(0, 0, width, height);
}

// void render_text(Shader &shader, string text, float x, float y, float z, float scale, glm::vec3 color)
// {
//     shader.use();
//     glUniform3f(glGetUniformLocation(shader.ID, "textColor"), color.x, color.y, color.z);
//     glActiveTexture(GL_TEXTURE0);
//     glBindVertexArray(VAO);

//     string::const_iterator c;
//     for (c = text.begin(); c != text.end(); c++)
//     {
//         Character ch = Characters[*c];

//         float xpos = x + ch.Bearing.x * scale;
//         float ypos = y - (ch.Size.y - ch.Bearing.y) * scale;
//         float zpos = z;

//         float w = ch.Size.x * scale;
//         float h = ch.Size.y * scale;
//         float vertices[6][5] = {{xpos, ypos + h, zpos, 0.0f, 0.0f}, {xpos, ypos, zpos, 0.0f, 1.0f},
//                                 {xpos + w, ypos, zpos, 1.0f, 1.0f}, {xpos, ypos + h, zpos, 0.0f, 0.0f},
//                                 {xpos + w, ypos, zpos, 1.0f, 1.0f}, {xpos + w, ypos + h, zpos, 1.0f, 0.0f}};
//         glBindTexture(GL_TEXTURE_2D, ch.TextureID);
//         glBindBuffer(GL_ARRAY_BUFFER, VBO);
//         glBufferSubData(GL_ARRAY_BUFFER, 0, sizeof(vertices), vertices);

//         glBindBuffer(GL_ARRAY_BUFFER, 0);
//         glDrawArrays(GL_TRIANGLES, 0, 6);
//         x += (ch.Advance >> 6) * scale;
//     }
//     glBindVertexArray(0);
//     glBindTexture(GL_TEXTURE_2D, 0);
// }

void mouse_callback(GLFWwindow *window, double xposIn, double yposIn)
{
    float xpos = static_cast<float>(xposIn);
    float ypos = static_cast<float>(yposIn);
    if (firstMouse)
    {
        lastX = xpos;
        lastY = ypos;
        firstMouse = false;
    }

    float xoffset = xpos - lastX;
    float yoffset = lastY - ypos;

    lastX = xpos;
    lastY = ypos;

    camera.ProcessMouseMovement(xoffset, yoffset);
}

unsigned int loadTexture(char const *path, bool gammaCorrection = true)
{
    unsigned int textureID;
    glGenTextures(1, &textureID);

    int width, height, nrComponents;
    unsigned char *data = stbi_load(path, &width, &height, &nrComponents, 0);
    if (data)
    {
        GLenum internalFormat;
        GLenum dataFormat;
        if (nrComponents == 1)
        {
            internalFormat = dataFormat = GL_RED;
        }
        else if (nrComponents == 3)
        {
            internalFormat = gammaCorrection ? GL_SRGB : GL_RGB;
            dataFormat = GL_RGB;
        }
        else if (nrComponents == 4)
        {
            internalFormat = gammaCorrection ? GL_SRGB_ALPHA : GL_RGBA;
            dataFormat = GL_RGBA;
        }

        glBindTexture(GL_TEXTURE_2D, textureID);
        glTexImage2D(GL_TEXTURE_2D, 0, internalFormat, width, height, 0, dataFormat, GL_UNSIGNED_BYTE, data);
        glGenerateMipmap(GL_TEXTURE_2D);

        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_REPEAT);
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_REPEAT);
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR_MIPMAP_LINEAR);
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR);

        stbi_image_free(data);
    }
    else
    {
        cout << "Texture failed to load at path: " << path << endl;
        stbi_image_free(data);
    }

    return textureID;
}

const string drawing_file = "../../bulletserver/cmd/bulletserver/GL_DRAWING.txt";
const string buffer_file = "../../bulletserver/cmd/bulletserver/GL_BUFFER.txt";
constexpr int MAX_LINE = 30;
static int readlinecount = 0;

vector<shared_ptr<Bullet>> buffer_drawing_file_and_read()
{
    vector<shared_ptr<Bullet>> v;
    ifstream df(drawing_file, ios::in);
    fstream bf(buffer_file, ios::in | ios::out | ios::trunc);
    if (!df)
    {
        cerr << "Error reading drawing_file!\n";
        return v;
    }
    if (!bf)
    {
        cerr << "Error reading buffer_file!\n";
        return v;
    }

    string line;
    for (int i = 0; i < readlinecount; i++)
        getline(df, line);

    int thistimeread = 0;
    while (thistimeread < MAX_LINE && getline(df, line))
    {
        bf << line;
        thistimeread++;
    }

    readlinecount += thistimeread;

    Bullet tmp;
    for (int i = 0; i < thistimeread; i++)
    {
        bf >> tmp.DurationSecs >> tmp.Size;
        for (int j = 0; j < 4; j++)
            bf >> tmp.Color[j];
        for (int j = 0; j < 3; j++)
            bf >> tmp.Position[j];
        bf >> tmp.s;
        shared_ptr<Bullet> newbullet = make_shared<Bullet>(tmp);
        v.emplace_back(newbullet);
    }
    return v;
}

static unsigned bulletVAO, bulletVBO;
static float bulletvertices[9];
void draw_a_bullet(const Bullet &bullet)
{
    glGenVertexArrays(1, &bulletVAO);
    glGenBuffers(1, &bulletVBO);

    bulletvertices[0] = bullet.DurationSecs;
    bulletvertices[1] = bullet.Size;
    for (int j = 2; j < 6; j++)
        bulletvertices[j] = bullet.Color[j - 2];
    for (int j = 6; j < 9; j++)
        bulletvertices[j] = bullet.Position[j - 6];

    glBindBuffer(GL_ARRAY_BUFFER, bulletVBO);
    glBufferData(GL_ARRAY_BUFFER, sizeof(bulletvertices), bulletvertices, GL_STATIC_DRAW);

    glBindVertexArray(bulletVAO);
    glEnableVertexAttribArray(0);
    glVertexAttribPointer(0, 1, GL_FLOAT, GL_FALSE, sizeof(float) * 9, (void *)0);

    glEnableVertexAttribArray(1);
    glVertexAttribPointer(1, 1, GL_FLOAT, GL_FALSE, sizeof(float) * 9, (void *)(sizeof(float)));

    glEnableVertexAttribArray(2);
    glVertexAttribPointer(2, 4, GL_FLOAT, GL_FALSE, sizeof(float) * 9, (void *)(sizeof(float) * 2));

    glEnableVertexAttribArray(3);
    glVertexAttribPointer(3, 4, GL_FLOAT, GL_FALSE, sizeof(float) * 9, (void *)(sizeof(float) * 6));

    glBindVertexArray(bulletVAO);
    glDrawArrays(GL_TRIANGLES, 0, 1);
    glBindVertexArray(0);
}

constexpr float radius = 1.0f;
constexpr int slices = 20;
constexpr int stacks = 20;

float standardcircle[stacks * slices * 6];

void generateSphere()
{
    float x, y, z, stackAngle, sinStack, cosStack, sliceAngle, sinSlice, cosSlice;
    int idx = 0;
    for (int i = 0; i < stacks; ++i)
    {
        stackAngle = M_PI * i / stacks;
        sinStack = sinf(stackAngle);
        cosStack = cosf(stackAngle);

        for (int j = 0; j < slices; ++j)
        {
            sliceAngle = 2 * M_PI * j / slices;
            sinSlice = sinf(sliceAngle);
            cosSlice = cosf(sliceAngle);

            x = radius * sinStack * cosSlice;
            y = radius * sinStack * sinSlice;
            z = radius * cosStack;

            standardcircle[idx++] = x;
            standardcircle[idx++] = y;
            standardcircle[idx++] = z;
            standardcircle[idx++] = x / radius;
            standardcircle[idx++] = y / radius;
            standardcircle[idx++] = z / radius;
        }
    }
}
