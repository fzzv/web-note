import express from 'express';
import usersRouter from './routes/users';
import postsRouter from './routes/posts';

const app = express();
app.use(express.json());

app.use('/users', usersRouter);
app.use('/posts', postsRouter);
app.listen(3000, () => {
  console.log('🚀 Server is running on http://localhost:3000');
});
