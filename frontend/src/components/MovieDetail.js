import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Link, useParams, useNavigate } from 'react-router-dom';

const MovieDetail = () => {
  const { id } = useParams();
  const [movie, setMovie] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    axios.get(`http://localhost:8000/movies/${id}`)
      .then(response => {
        setMovie(response.data);
      })
      .catch(error => {
        console.error('There was an error fetching the movie!', error);
      });
  }, [id]);

  const deleteMovie = () => {
    axios.delete(`http://localhost:8000/movies/${id}`)
      .then(() => {
        navigate('/');
      })
      .catch(error => {
        console.error('There was an error deleting the movie!', error);
      });
  };

  if (!movie) return <div>Loading...</div>;

  return (
    <div>
      <h2>{movie.title}</h2>
      <p>ISBN: {movie.isbn}</p>
      <p>Director: {movie.director.first_name} {movie.director.last_name}</p>
      <button onClick={() => navigate(`/edit-movie/${id}`)}>Edit</button>
      <button onClick={deleteMovie}>Delete</button>
      <br />
      <Link to="/">Back to Movie List</Link>
    </div>
  );
};

export default MovieDetail;
