<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Report Bun</title>
  </head>
  <body>
    <form action="/bunReport" method="post">
      Size: <input type="number" name="size" min="1" max="10">
      Description: <input type="text" name="description">
      <input type="submit" value="Report">
    </form>
  </body>
</html>
