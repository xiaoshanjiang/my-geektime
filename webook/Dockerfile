# 基础镜像
FROM ubuntu:22.04
# 把编译后的打包进来这个镜像，放到工作目录 /app。你随便换
COPY ./webook-app /app/webook-app
WORKDIR /app
# ENTRYPOINT和CMD的区别: https://docs.docker.com/engine/reference/builder/#entrypoint
RUN chmod a+x ./webook-app
CMD ["./webook-app"]
