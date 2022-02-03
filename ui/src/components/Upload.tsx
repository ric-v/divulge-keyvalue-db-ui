import React, { FormEvent, useState } from 'react';
import { Button } from "@mui/material";
import { useSnackbar } from 'notistack';
import http from '../services/axios-common';
import { CircularProgressbar, buildStyles } from 'react-circular-progressbar';
import "react-circular-progressbar/dist/styles.css";

const Upload = () => {

  const { enqueueSnackbar, closeSnackbar } = useSnackbar();
  const [file, setFile] = useState<string | Blob | File>('');

  const [uploadProgress, updateUploadProgress] = useState(0);
  const [uploadStatus, setUploadStatus] = useState(false);
  const [uploading, setUploading] = useState(false);

  const handleFileUpload = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    setUploading(true);
    const formData = new FormData();
    formData.append('file', file);

    http.post('/api/v1/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (ev: ProgressEvent) => {
        const progress = ev.loaded / ev.total * 100;
        updateUploadProgress(Math.round(progress));
      },
    })
      .then((resp) => {
        // our mocked response will always return true
        // in practice, you would want to use the actual response object
        setUploadStatus(true);
        setUploading(false);
        console.log(resp);

        enqueueSnackbar('success', {
          key: 'success',
          variant: 'success',
          draggable: true,
          onClick: () => {
            closeSnackbar('success');
          }
        });
      })
      .catch((err) => console.error(err));
  };

  return (
    <form onSubmit={handleFileUpload}>
      <h1>Select your files</h1>
      <input
        type='file'
        name='file-choose'
        onChange={(e) => {
          e.preventDefault();
          if (e.target.files && e.target.files.length > 0) {
            setFile(e.target.files[0])
          }
        }}
      />
      <Button
        variant='contained'
        color="secondary"
        type='submit'
      >
        Upload
      </Button>
      {(uploading)
        ?
        <div className="progress-bar-container">
          <CircularProgressbar
            value={uploadProgress}
            text={`${uploadProgress}% uploaded`}
            styles={buildStyles({
              textSize: '10px',
              pathColor: 'teal',
            })}
          />
        </div>
        : null
      }
    </form >
  );
};

export default Upload;

    // // GET the file from API server
    // fetch('http://localhost:8080/api/v1/new?file='  fileName  '&dbtype=boltdb', {
    //   method: 'POST',
    // })
    //   .then(response => response.json())
    //   .then(data => {
    //     console.log('success', data);
    //     // store the access key to local storage
    //     localStorage.setItem('access_key', data.dbkey);
    //     localStorage.setItem('file_name', data.filename);
    //     localStorage.setItem('db_type', data.dbtype);
    //     console.log('success', data);

    //     enqueueSnackbar('success', {
    //       key: 'success',
    //       variant: 'success',
    //       draggable: true,
    //       onClick: () => {
    //         closeSnackbar('success');
    //       }
    //     });
    //   })
    //   .catch(err => {
    //     console.log('error', err);
    //   });