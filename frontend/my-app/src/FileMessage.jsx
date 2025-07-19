// import React from "react";
// import "./FileMessage.css"; // import CSS file

// const FileMessage = ({ fileName, fileSize, downloadUrl }) => {
//   return (
//     <div className="file-message-container">
//       <div className="file-message-icon">📄</div>
//       <div className="file-message-info">
//         <div className="file-message-name">{fileName}</div>
//         <div className="file-message-size">{fileSize}</div>
//       </div>
//       <a href={downloadUrl} download className="file-message-button">
//         Download
//       </a>
//     </div>
//   );
// };

// export default FileMessage;

import "./FileMessage.css";


function FileMessage({ fileName, fileSize, downloadUrl }) {

    function triggerDownload(url, fileName) {
  const a = document.createElement("a");
  a.href = url;
  a.download = fileName; // Force download instead of navigating
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
}
  return (
    <div className="file-message">
      <div className="file-details">
        <strong>{fileName}</strong> <span>{fileSize}</span>
      </div>
      <a href={downloadUrl} target="_blank" rel="noopener noreferrer">
        <button className="download-btn" onClick={() => triggerDownload(downloadUrl, "Smaple File")} >Download</button>
      </a>
    </div>
  );
}

export default FileMessage;
