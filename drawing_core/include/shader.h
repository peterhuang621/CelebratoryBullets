#ifndef SHADER_H
#define SHADER_H

#include <OpenGL/gl3.h>
#include <fstream>
#include <glm/glm.hpp>
#include <iostream>
#include <sstream>
using namespace std;

class Shader
{
  public:
    unsigned int ID;
    Shader(const char *vertexPath, const char *fragmentPath, const char *geomertryPath);
    Shader(const string vertexPath, const string fragmentPath, const string geomertryPath);
    Shader(const char *vertexPath, const char *fragmentPath);
    Shader(const string vertexPath, const string fragmentPath);

    void use();
    void setBool(const string &name, bool value) const;
    void setInt(const string &name, int value) const;
    void setFloat(const string &name, float value) const;
    void setVec2(const string &name, const glm::vec2 &value) const;
    void setVec2(const std::string &name, float x, float y) const;
    void setVec3(const std::string &name, const glm::vec3 &value) const;
    void setVec3(const std::string &name, float x, float y, float z) const;
    void setVec4(const std::string &name, const glm::vec4 &value) const;
    void setVec4(const std::string &name, float x, float y, float z, float w) const;
    void setMat2(const std::string &name, const glm::mat2 &mat) const;
    void setMat3(const std::string &name, const glm::mat3 &mat) const;
    void setMat4(const std::string &name, const glm::mat4 &mat) const;

  private:
    void checkCompileErrors(unsigned int shader, string type);
};

#endif