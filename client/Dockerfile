# Stage 1
FROM node:16 as react-build
WORKDIR /app
COPY ./client ./
RUN npm install
RUN npm run build

# Stage 2 - the production environment
FROM nginx:1.20-alpine
COPY ./client/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=react-build /app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]