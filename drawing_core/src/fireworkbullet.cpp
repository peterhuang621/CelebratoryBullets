#include "fireworkbullet.h"

int main()
{
    glfwInit();
    glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 3);
    glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 3);
    glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE);

#ifdef __APPLE__
    glfwWindowHint(GLFW_OPENGL_FORWARD_COMPAT, GL_TRUE);
#endif

    GLFWwindow *window = glfwCreateWindow(SCR_WIDTH, SCR_HEIGHT, "CelebratoryBullets", NULL, NULL);
    if (window == NULL)
    {
        cout << "Failed to create GLFW window" << endl;
        glfwTerminate();
        return -1;
    }
    glfwMakeContextCurrent(window);
    glfwSetFramebufferSizeCallback(window, framebuffer_size_callback);
    glfwSetCursorPosCallback(window, mouse_callback);
    // glEnable(GL_PROGRAM_POINT_SIZE);
    glEnable(GL_DEPTH);
    Shader bulletshader("/Users/peterhuang98/test_code/Go/CelebratoryBullets/drawing_core/src/bullet.vs",
                        "/Users/peterhuang98/test_code/Go/CelebratoryBullets/drawing_core/src/bullet.fs",
                        "/Users/peterhuang98/test_code/Go/CelebratoryBullets/drawing_core/src/bullet.gs");
    generateSphere();
    unsigned vao, vbo;
    glGenVertexArrays(1, &vao);
    glBindVertexArray(vao);

    glGenBuffers(1, &vbo);
    glBindBuffer(GL_ARRAY_BUFFER, vbo);
    glBufferData(GL_ARRAY_BUFFER, sizeof(standardcircle), standardcircle, GL_STATIC_DRAW);

    glEnableVertexAttribArray(0);
    glVertexAttribPointer(0, 3, GL_FLOAT, GL_FALSE, 6 * sizeof(float), (void *)0);
    glEnableVertexAttribArray(1);
    glVertexAttribPointer(1, 3, GL_FLOAT, GL_FALSE, 6 * sizeof(float), (void *)(sizeof(float) * 3));
    glBindVertexArray(0);

    const glm::mat4 projection =
        glm::perspective(glm::radians(45.0f), (float)SCR_WIDTH / (float)SCR_HEIGHT, 0.1f, 1000.0f);
    glm::mat4 view = camera.GetViewMatrix();
    glm::mat4 model = glm::mat4(1.0f);
    float currentFrame, dis;

    while (!glfwWindowShouldClose(window))
    {
        currentFrame = static_cast<float>(glfwGetTime());
        deltaTime = currentFrame - lastFrame;
        lastFrame = currentFrame;
        processInput(window);

        glClearColor(0.2f, 0.3f, 0.3f, 1.0f);
        glClear(GL_COLOR_BUFFER_BIT);

        bulletshader.use();
        bulletshader.setMat4("projection", projection);
        view = camera.GetViewMatrix();
        bulletshader.setMat4("view", view);
        dis = sinf(currentFrame);
        bulletshader.setFloat("dis", dis);
        bulletshader.setMat4("model", model);
        glBindVertexArray(vao);
        glDrawArrays(GL_POINTS, 0, stacks * slices);

        glfwSwapBuffers(window);
        glfwPollEvents();
    }

    glfwTerminate();
    return 0;
}
