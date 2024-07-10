import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate, useParams } from 'react-router-dom';

const MovieForm = () => {
  const [movie, setMovie] = useState({ title: '', isbn: '', director: { first_name: '', last_name: '' } });
  const navigate = useNavigate();
  const { id } = useParams();

  useEffect(() => {
    if (id) {
      axios.get(`http://localhost:8000/movies/${id}`)
        .then(response => {
          setMovie(response.data);
        })
        .catch(error => {
          console.error('There was an error fetching the movie!', error);
        });
    }
  }, [id]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    if (name === 'first_name' || name === 'last_name') {
      setMovie({ ...movie, director: { ...movie.director, [name]: value } });
    } else {
      setMovie({ ...movie, [name]: value });
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (id) {
      axios.put(`http://localhost:8000/movies/${id}`, movie)
        .then(() => {
          navigate('/');
        })
        .catch(error => {
          console.error('There was an error updating the movie!', error);
        });
    } else {
      axios.post('http://localhost:8000/movies', movie)
        .then(() => {
          navigate('/');
        })
        .catch(error => {
          console.error('There was an error creating the movie!', error);
        });
    }
  };

  return (
    <div>
      <h2>{id ? 'Edit Movie' : 'Add Movie'}</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Title:</label>
          <input type="text" name="title" value={movie.title} onChange={handleChange} />
        </div>
        <div>
          <label>ISBN:</label>
          <input type="text" name="isbn" value={movie.isbn} onChange={handleChange} />
        </div>
        <div>
          <label>Director First Name:</label>
          <input type="text" name="first_name" value={movie.director.first_name} onChange={handleChange} />
        </div>
        <div>
          <label>Director Last Name:</label>
          <input type="text" name="last_name" value={movie.director.last_name} onChange={handleChange} />
        </div>
        <button type="submit">{id ? 'Update' : 'Create'}</button>
      </form>
    </div>
  );
};

export default MovieForm;
