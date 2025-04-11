// #version 330 core
// layout (points) in;
// layout (triangle_strip, max_vertices = 5) out;

// in VS_OUT {
//     vec3 color;
// } gs_in[];

// out vec3 fColor;

// uniform mat4 projection;
// uniform mat4 view;
// uniform mat4 model;
// uniform float dis;

// void build_house(vec4 position)
// {    
//     fColor = gs_in[0].color; // gs_in[0] since there's only one input vertex
//     gl_Position = projection*view*model *(position + vec4(-0.02, -0.02, 0.0, 0.0)); // 1:bottom-left   
//     EmitVertex();   
//     gl_Position = projection*view*model *(position + vec4( 0.02, -0.02, 0.0, 0.0)); // 2:bottom-right
//     EmitVertex();
//     gl_Position = projection*view*model *(position + vec4(-0.02,  0.02, 0.0, 0.0)); // 3:top-left
//     EmitVertex();
//     // gl_Position = projection*view*model *(position + vec4( 0.2,  0.2, 0.0, 0.0)); // 4:top-right
//     // EmitVertex();
//     // gl_Position = projection*view*model *(position + vec4( 0.0,  0.4, 0.0, 0.0)); // 5:top
//     // EmitVertex();
//     EndPrimitive();
// }

// void main() {    
//     build_house(gl_in[0].gl_Position);
// }

#version 330 core
layout(points) in;
layout(triangle_strip, max_vertices = 3) out;

in vec3 vPos[];
in vec3 vNormal[];

out vec3 fNormal;

uniform float start,end,time;
uniform mat4 projection,view,model;
vec3 down;

vec3 pointmove(){
    vec3 ret=vec3(0.0,0.0,0.0);
    float ratio = (time-start) / (end-start);
    ret += (down*(1-vPos[0][1])*2 + vNormal[0])*ratio;
    return ret;
}

void main() {
    down = vec3(0.0, -1.0, 0.0);
    vec3 center = vPos[0];
    vec3 normal = normalize(vNormal[0]);

    // 生成一个朝法线方向稍微展开的小三角形
    float size = 0.05; // 三角形大小

    vec3 tangent = normalize(cross(normal, vec3(0.0, 1.0, 0.0)));
    if (length(tangent) < 0.01)
        tangent = normalize(cross(normal, vec3(1.0, 0.0, 0.0)));

    vec3 bitangent = normalize(cross(normal, tangent));
    
    vec3 moves=pointmove();
    vec3 p1 = center + normal * size + moves;
    vec3 p2 = center + tangent * size * 0.5 - normal * size * 0.2 + moves;
    vec3 p3 = center + bitangent * size * 0.5 - normal * size * 0.2 + moves;

    gl_Position = projection * view * vec4(p1, 1.0);
    fNormal = normal;
    EmitVertex();

    gl_Position = projection * view * vec4(p2, 1.0);
    fNormal = normal;
    EmitVertex();

    gl_Position = projection * view * vec4(p3, 1.0);
    fNormal = normal;
    EmitVertex();

    EndPrimitive();
}
