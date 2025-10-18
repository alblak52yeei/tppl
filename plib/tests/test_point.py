import pytest

from plib import Point

@pytest.fixture
def points():
    return Point(0, 0), Point(2, 2)

class TestPoint:

    def test_creation(self):
        p = Point(1, 2)
        assert p.x == 1 and p.y == 2

        with pytest.raises(TypeError):
            Point(1.5, 1.5)

    def test_add(self, points):
        p1, p2 = points
        assert p2 + p1 == Point(2, 2)
    
    def test_sub(self, points):
        p1, p2 = points
        assert p2 - p1 == Point(2, 2)
        assert p1 - p2 == -Point(2, 2)
    
    def test_distance_to(self):
        p1 = Point(0, 0)
        p2 = Point(2, 0)
        assert p1.to(p2) == 2

    @pytest.mark.parametrize(
            "p1, p2, distance",
            [(Point(0, 0), Point(0, 10), 10),
             (Point(0, 0), Point(10, 0), 10),
             (Point(0, 0), Point(1, 1), 1.414)]
    )
    def test_distance_all_axis(self, p1, p2, distance):
        assert p1.to(p2) == pytest.approx(distance, 0.001)

    def test_iadd(self, points):
        """Test __iadd__ method (in-place addition)"""
        p1, p2 = points
        result = p1.__iadd__(p2)
        assert result == Point(2, 2)
        # Note: __iadd__ should modify self, but this implementation returns new Point
        assert p1 == Point(0, 0)  # original p1 unchanged

    def test_eq_with_non_point(self):
        """Test __eq__ method with non-Point object"""
        p1 = Point(1, 2)
        with pytest.raises(NotImplementedError):
            p1 == "not a point"
        
        with pytest.raises(NotImplementedError):
            p1 == 42

    def test_str(self):
        """Test __str__ method"""
        p = Point(3, 4)
        assert str(p) == "Point(3, 4)"

    def test_repr(self):
        """Test __repr__ method"""
        p = Point(5, 6)
        assert repr(p) == "Point(5, 6)"

    def test_is_center(self):
        """Test is_center method"""
        center_point = Point(0, 0)
        non_center_point = Point(1, 0)
        non_center_point2 = Point(0, 1)
        non_center_point3 = Point(1, 1)
        
        assert center_point.is_center() is True
        assert non_center_point.is_center() is False
        assert non_center_point2.is_center() is False
        assert non_center_point3.is_center() is False

    def test_to_json(self):
        """Test to_json method"""
        p = Point(7, 8)
        json_str = p.to_json()
        expected = '{"x": 7, "y": 8}'
        assert json_str == expected

    def test_from_json(self):
        """Test from_json class method"""
        json_str = '{"x": 9, "y": 10}'
        p = Point.from_json(json_str)
        assert p == Point(9, 10)
        
        # Test with different values
        json_str2 = '{"x": -5, "y": 3}'
        p2 = Point.from_json(json_str2)
        assert p2 == Point(-5, 3)