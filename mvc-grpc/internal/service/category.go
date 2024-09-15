package service

import (
	"context"
	"io"

	"github.com/markuscandido/mvc-grpc/internal/database"
	"github.com/markuscandido/mvc-grpc/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB: categoryDB,
	}
}

func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	category, err := c.CategoryDB.Create(in.Name, in.Description)
	if err != nil {
		return nil, err
	}
	createCategoryResponse := &pb.CreateCategoryResponse{
		Category: &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		},
	}

	return createCategoryResponse, nil
}

func (c *CategoryService) ListCategory(ctx context.Context, in *pb.ListCategoryRequest) (*pb.ListCategoryResponse, error) {
	categories, err := c.CategoryDB.FindAll()
	if err != nil {
		return nil, err
	}

	var categoriesResponse []*pb.Category
	for _, category := range categories {
		categoryResponse := &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		}

		categoriesResponse = append(categoriesResponse, categoryResponse)
	}

	return &pb.ListCategoryResponse{
		Categories: categoriesResponse,
	}, nil
}

func (c *CategoryService) GetCategory(ctx context.Context, in *pb.GetCategoryRequest) (*pb.GetCategoryResponse, error) {
	category, err := c.CategoryDB.FindByID(in.Id)
	if err != nil {
		return nil, err
	}
	categoryResponse := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return &pb.GetCategoryResponse{
		Category: categoryResponse,
	}, nil
}

func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	response := &pb.CreateCategoryStreamResponse{}
	for {
		categoryRequest, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(response)
		}
		if err != nil {
			return err
		}

		categoryDB, err := c.CategoryDB.Create(categoryRequest.Name, categoryRequest.Description)
		if err != nil {
			return err
		}

		response.Categories = append(response.Categories, &pb.Category{
			Id:          categoryDB.ID,
			Name:        categoryDB.Name,
			Description: categoryDB.Description,
		})
	}
}

func (c *CategoryService) CreateCategoryStreamBidirectional(stream pb.CategoryService_CreateCategoryStreamBidirectionalServer) error {
	for {
		categoryRequest, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		categoryDB, err := c.CategoryDB.Create(categoryRequest.Name, categoryRequest.Description)
		if err != nil {
			return err
		}

		err = stream.Send(&pb.CreateCategoryStreamBidirectionalResponse{
			Category: &pb.Category{
				Id:          categoryDB.ID,
				Name:        categoryDB.Name,
				Description: categoryDB.Description,
			},
		})
		if err != nil {
			return err
		}
	}
}
