
const fs = require('fs');
const path = require('path');

// https://www.codota.com/code/javascript/functions/fs/rmdirSync
// WHY? rmdirSync() does not return on time... 

function rmdir(dir) {
 if (!fs.existsSync(dir)) {
  return null;
 }
 fs.readdirSync(dir).forEach(f => {
  let pathname = path.join(dir, f);
  if (!fs.existsSync(pathname)) {
   return fs.unlinkSync(pathname);
  }
  if (fs.statSync(pathname).isDirectory()) {
   return rmdir(pathname);
  } else {
   return fs.unlinkSync(pathname);
  }
 });
 return fs.rmdirSync(dir);
}

module.exports = {
  rmdir
}

